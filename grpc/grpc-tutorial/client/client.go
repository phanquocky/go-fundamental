package main

import (
	"context"
	"fmt"
	routeguide "go-fundamental/grpc/grpc-tutorial/code-gen"
	"io"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	cred := insecure.NewCredentials()
	client, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatal("cannot create client to localhost:8080", err.Error())
	}

	routeGuideService := routeguide.NewRouteGuideClient(client)

	// test GetFeature function
	feature, err := routeGuideService.GetFeature(context.Background(), &routeguide.Point{Latitude: 409146138, Longitude: -746188906})
	if err != nil {
		log.Fatal("cannot call getfeature:", err)
	}

	fmt.Println("Getfeatures: ", feature)

	featureMiss, err := routeGuideService.GetFeature(context.Background(), &routeguide.Point{Latitude: 1, Longitude: 0})
	if err != nil {
		log.Fatal("cannot call getfeature miss:", err)
	}

	fmt.Println("Getfeatures: ", featureMiss)

	// test ListFeature function
	listFeaturesStream, err := routeGuideService.ListFeatures(context.Background(), &routeguide.Rectangle{
		Lo: &routeguide.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &routeguide.Point{Latitude: 420000000, Longitude: -730000000},
	})
	if err != nil {
		log.Fatal("cannot call ListFeatures: ", err)
	}

	for {
		feature, err := listFeaturesStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot receive feature from client from ListFeature func", err)
		}

		fmt.Println("[ListFeatures] feature: ", feature)
	}

	// test recodeRoute function
	// Create a random number of random points
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pointCount := int(r.Int31n(100)) + 2 // Traverse at least two points
	var points []*routeguide.Point
	for i := 0; i < pointCount; i++ {
		points = append(points, randomPoint(r))
	}
	log.Printf("Traversing %d points.", len(points))
	clientStream, err := routeGuideService.RecordRoute(context.Background())
	if err != nil {
		log.Fatalf("%v.RecordRoute(_) = _, %v", client, err)
	}

	for _, point := range points {
		if err := clientStream.Send(point); err != nil {
			log.Fatalf("%v.Send(%v) = %v", clientStream, point, err)
		}
	}

	reply, err := clientStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", clientStream, err, nil)
	}
	log.Printf("Route summary: %v", reply)

	// test RouteChat func
	notes := []*routeguide.RouteNote{
		{Location: &routeguide.Point{Latitude: 0, Longitude: 1}, Message: "First message"},
		{Location: &routeguide.Point{Latitude: 0, Longitude: 2}, Message: "Second message"},
		{Location: &routeguide.Point{Latitude: 0, Longitude: 3}, Message: "Third message"},
		{Location: &routeguide.Point{Latitude: 0, Longitude: 1}, Message: "Fourth message"},
		{Location: &routeguide.Point{Latitude: 0, Longitude: 2}, Message: "Fifth message"},
		{Location: &routeguide.Point{Latitude: 0, Longitude: 3}, Message: "Sixth message"},
	}

	RouteChatStream, err := routeGuideService.RouteChat(context.Background())
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := RouteChatStream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			log.Printf("[RouteChat] Got message %s at point(%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
		}
	}()
	for _, note := range notes {
		if err := RouteChatStream.Send(note); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}
	}
	RouteChatStream.CloseSend()
	<-waitc
}

func randomPoint(r *rand.Rand) *routeguide.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &routeguide.Point{Latitude: lat, Longitude: long}
}
