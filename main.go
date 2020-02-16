package main

import (
    "fmt"
    "os"
    "strings"
    "bufio"
    "strconv"
    "sort"
    "math"
    "math/rand"
    "time"
)

const (
    FileName = "gifts.txt"
    MaximumWeight = 10000 * 1000
    StartLat = 68.073611
    StartLong = 29.315278
    EarthRadius = 6378.0
    InitialTemperature = 1000.0
    CoolingRate = 0.003
)

type Gift struct {
    childId int
    latitude float64
    longitude float64
    grams int
}

func parseInt(childId string, grams string) (int, int) {
    result1, err := strconv.ParseUint(childId, 10, 32)
    result2, err2 := strconv.ParseUint(grams, 10, 32)
    if err != nil || err2 != nil {
        return -1, -1
    }
    return int(result1), int(result2)
}

func parseFloat(latitude string, longitude string) (float64, float64) {
    result1, err := strconv.ParseFloat(latitude, 64)
    result2, err2 := strconv.ParseFloat(longitude, 64)
    if err != nil || err2 != nil {
        return -1, -1
    }
    return result1, result2
}

func parseGiftFile(file string) ([]Gift, error) {
    pwd, _ := os.Getwd()
    fh, fileErr := os.Open(pwd + "/" + FileName)
    if fileErr != nil {
        fmt.Printf("Error in opening file: %v\n", fileErr)
        return nil, fileErr
    }
    defer fh.Close()

    var gifts []Gift
    scanner := bufio.NewScanner(fh)
    for scanner.Scan() {
        gift := strings.Split(scanner.Text(), ";")
        childId, grams := parseInt(gift[0], gift[3])
        latitude, longitude := parseFloat(gift[1], gift[2])
        if childId == -1 || grams == -1 || latitude == -1 || longitude == -1 { continue }
        gifts = append(gifts, Gift{childId, latitude, longitude, grams})
    }
    return gifts, nil
}

func toRadians(degree float64) float64 {
    return (degree * (math.Pi / 180.0))
}

func distanceBetweenTwoPoints(lat1 float64, long1 float64, lat2 float64, long2 float64) float64 {
    lat1, long1, lat2, long2 = toRadians(lat1), toRadians(long1), toRadians(lat2), toRadians(long2)
    dlong, dlat := long2 - long1, lat2 - lat1
    result := math.Pow(math.Sin(dlat / 2), 2.0) + math.Cos(lat1) * math.Cos(lat2) * math.Pow(math.Sin(dlong / 2), 2.0)
    result = 2 * math.Atan2(math.Sqrt(result), math.Sqrt(1-result))
    return result * EarthRadius
}

// Calculate the distance using radius, lat and long positions
func calculateDistance(trips []Gift) float64 {
    distance := 0.0
    for i := 0; i < len(trips); i++ {
        lat1, long1, lat2, long2 := 0.0, 0.0, 0.0, 0.0
        if i == 0 || i == len(trips) - 1 {
            if i == 0 {
                lat1, long1, lat2, long2 = StartLat, StartLong, trips[i].latitude, trips[i].longitude
            } else if i == len(trips) - 1 {
                lat1, long1, lat2, long2 = trips[i].latitude, trips[i].longitude, StartLat, StartLong
            }
        } else {
            lat1, long1, lat2, long2 = trips[i].latitude, trips[i].longitude, trips[i+1].latitude, trips[i+1].longitude
        }
        distance = distance + distanceBetweenTwoPoints(lat1, long1, lat2, long2)
    }
    return distance
}


/*
This algorithm uses the greedy approach to solve this problem.
Sort the gifts by their weight in grams. Pick one of them each from each sides
until you hit the maximum weight, calculate the distance and repeat until there
are no gifts
*/
func solveGreedy(gifts []Gift) float64 {
    sort.Slice(gifts, func(i int, j int) bool {
        return gifts[i].grams < gifts[j].grams
    })
    i, j := 0, len(gifts) - 1
    var trips []Gift
    var currentWeight int
    currentWeight = 0
    totalDistance := 0.0
    for i <= j {
        foundGiftThatCanFit := false
        if currentWeight + gifts[i].grams <= MaximumWeight {
            currentWeight = currentWeight + gifts[i].grams
            trips = append(trips, gifts[i])
            i = i + 1
            foundGiftThatCanFit = true
        }
        if currentWeight + gifts[j].grams <= MaximumWeight {
            currentWeight = currentWeight + gifts[j].grams
            trips = append(trips, gifts[j])
            j = j - 1
            foundGiftThatCanFit = true
        }
        if foundGiftThatCanFit == false {
            distancePerTrip := calculateDistance(trips)
            totalDistance = totalDistance + distancePerTrip
            currentWeight = 0
            trips = nil
        }
    }
    return totalDistance * 1000
}

/*
Simulated Annealing. For reference, follow these resources

1) https://en.wikipedia.org/wiki/Simulated_annealing
2) http://www.theprojectspot.com/tutorial-post/simulated-annealing-algorithm-for-beginners/6
*/
func acceptanceProbability(energy float64, newEnergy float64, temperature float64) float64 {
    if newEnergy < energy { return 1.0 }
    return math.Exp((energy - newEnergy) / temperature)
}

func randomShuffle(gifts []Gift) []Gift {
    rand.Seed(time.Now().UnixNano())
    rand.Shuffle(len(gifts), func(i, j int) { gifts[i], gifts[j] = gifts[j], gifts[i] })
    return gifts
}

func getNeighbour(gifts []Gift) []Gift {
    rand.Seed(time.Now().UnixNano())
    randomPos1 := rand.Intn(len(gifts))
    randomPos2 := rand.Intn(len(gifts))
    gifts[randomPos1], gifts[randomPos2] = gifts[randomPos2], gifts[randomPos1]
    return gifts
}

func getTotalTripDistance(gifts []Gift) float64 {
    totalDistance := 0.0
    curWeight := 0
    var trips []Gift
    for i := 0; i < len(gifts); i = i+1 {
        if curWeight + gifts[i].grams < MaximumWeight {
            curWeight = curWeight + gifts[i].grams
            trips = append(trips, gifts[i])
        } else {
            totalDistance = totalDistance + calculateDistance(trips)
            trips = nil
            curWeight = 0.0
        }
    }
    totalDistance = totalDistance + calculateDistance(trips)
    return totalDistance
}

func solveSimulatedAnnealing(gifts []Gift) float64 {
    currentSolution := randomShuffle(gifts)
    temperature := InitialTemperature
    bestDistance := 100000000000.0
    for temperature > 1 {
        neighborSolution := getNeighbour(gifts)
        curEnergy, newEnergy := getTotalTripDistance(currentSolution), getTotalTripDistance(neighborSolution)
        if acceptanceProbability(curEnergy, newEnergy, temperature) > rand.Float64() {
            currentSolution = neighborSolution
        }
        if getTotalTripDistance(currentSolution) < bestDistance {
            bestDistance = getTotalTripDistance(currentSolution)
        }
        temperature = temperature * (1 - CoolingRate)
    }
    return bestDistance * 1000
}

func main() {
    gifts, err := parseGiftFile(FileName)
    if err != nil {
        fmt.Printf("Error in running parseGiftFile method: %v\n", err)
        os.Exit(0)
    }
    result := solveGreedy(gifts)
    fmt.Printf("total distance via greedy approach: %d\n", int(result))
    result = solveSimulatedAnnealing(gifts)
    fmt.Printf("total distance via simulated annealing: %d\n", int(result))
}
