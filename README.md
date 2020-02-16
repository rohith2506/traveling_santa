## Problem Statement

At your disposal, you have Santa’s most highly guarded trade secret, the Nice List. The coveted List details which present shall be given to which child. Here’s the kicker: Santa’s sleigh can only carry 10 metric tons (10,000kg) at a time, and for each trip Santa makes, you'll need to tell him which items to pack.

You can download the list here. In the file, you'll find the wish of each child. Each row contains one wish. For privacy purposes, we've left out the names of the children: instead, their files include a numerical ID, their coordinates on Earth, and the weight of their present in grams. Your job is to find the most optimal routes to deliver all the presents. Santa starts at Korvatunturi, Finland (68.073611N 29.315278E) and, based on your list, flies directly from one coordinate to another until all presents are delivered, and then returns to Korvatunturi. The shorter the overall length of the trip, the less emissions there will be. For the purposes of this task, we assume Earth to be a sphere with radius of 6,378km.

[original problem link](https://traveling-santa.reaktor.com/)

## Solutions

As this is a double NP complete problem, there is no possible deterministic algorithm which can solve this. As gifts are in order of 10k, running through all combinations is impractical. We have to rely on heuristric and optimization methods.

1) Greedy Approach
2) Simulated Annealing

### Greedy Approach
Sort all the gifts given in order of their weights and then divided them into 10000 kg buckets. Mark each bucket as a trip and calculate distance

### Simulated Annealing
This is a classic optimization technique for NP complete problems. We start with a higher initial temperature, try to find the nearest neighbour, compare it with original solution and choose the best and keep repeating this loop until temperature is cooled enough.

For detailed overview, I would recommend this articles

http://www.theprojectspot.com/tutorial-post/simulated-annealing-algorithm-for-beginners/6
https://en.wikipedia.org/wiki/Simulated_annealing

## How to run
`go run main.go`

## Credits
Rohith Uppala < rohith.uppala369@gmail.com >
