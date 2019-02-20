//I Terrence Light (TE965355) affirm that this program is entirely my own work and that I have neither developed my code together with any another person, nor copied any code from any other person, nor permitted my code to be copied or otherwise used by any other person, nor have I copied, modified, or otherwise used programs created by others. I acknowledge that any violation of the above terms will be treated as academic dishonesty.

//Terrence Light te965355
//COP 4600 Operating Systems
//Project 1: Disk Scheduling Algorithms

package main

import (
	"bufio"
	"log"
	"fmt"
	"os"
	"strings"
	"strconv"
	"math"
)

//We need to use custom sorting functions because Eustis hasn't been updated...
const NADA int = -1

func DeepCopy(vals []int) []int {
   tmp := make([]int, len(vals))
   copy(tmp, vals)
   return tmp
}

func MergeSort(items []int) {

   if len(items) > 1 {
      mid := len(items) / 2
      left := DeepCopy(items[0:mid])
      right := DeepCopy(items[mid:])

      MergeSort(left)
      MergeSort(right)

      l := 0
      r := 0

      for i := 0; i < len(items); i++ {

         lval := NADA
         rval := NADA

         if l < len(left) {
            lval = left[l]
         }

         if r < len(right) {
            rval = right[r]
         }

         if (lval != NADA && rval != NADA && lval < rval) || rval == NADA {
            items[i] = lval
            l += 1
         } else if (lval != NADA && rval != NADA && lval >= rval) || lval == NADA {
            items[i] = rval
            r += 1
         }

      }
   }

}

func fcfs(requestArr []int, initCyl int, lowerCyl int, upperCyl int) {
	//First we need to do a data dump before procession the scheduling
	fmt.Printf("Seek algorithm: FCFS\n")
	fmt.Printf("\tLower cylinder: %5d\n", lowerCyl)
	fmt.Printf("\tUpper cylinder: %5d\n", upperCyl)
	fmt.Printf("\tInit cylinder:  %5d\n", initCyl)
	fmt.Printf("\tCylinder requests:\n")
	
	for i := 0; i < len(requestArr); i++ {
		fmt.Printf("\t\tCylinder %5d\n", requestArr[i])
	}
	
	//Track the number of cylinders we traverse
	traverse := 0
	moved := 0.0
	curCyl := initCyl
	
	//Since this is fcfs, we don't need to do any sorting or the like
	//Since we traverse down the whole slice, we don't need to mark the used bools for each request
	for i := 0; i < len(requestArr); i++ {
		//Make sure a request is valid, if not then skip it
		if((requestArr[i] > upperCyl) || (requestArr[i] < lowerCyl)) {
			fmt.Printf("ERROR(15):Request out of bounds:  req (%d) > upper (%d) or < lower (%d)\n", requestArr[i], upperCyl, lowerCyl)
			continue
		}
	
		//Service the next request in line and updating the number of cylinders traversed
		fmt.Printf("Servicing %5d\n", requestArr[i])
		moved = math.Abs(float64(curCyl - requestArr[i]))
		curCyl = requestArr[i]
		traverse = traverse + int(moved)
	}
	
	fmt.Printf("FCFS traversal count = %5d\n", traverse)
	
}

func sstf(requestArr []int, initCyl int, lowerCyl int, upperCyl int) {
	//First we need to do a data dump before procession the scheduling
	fmt.Printf("Seek algorithm: SSTF\n")
	fmt.Printf("\tLower cylinder: %5d\n", lowerCyl)
	fmt.Printf("\tUpper cylinder: %5d\n", upperCyl)
	fmt.Printf("\tInit cylinder:  %5d\n", initCyl)
	fmt.Printf("\tCylinder requests:\n")
	
	for i := 0; i < len(requestArr); i++ {
		fmt.Printf("\t\tCylinder %5d\n", requestArr[i])
	}
	
	//Variables for tracking seek time
	traverse := 0
	curCyl := initCyl
	initLength := len(requestArr)
	curShortest := math.MaxInt32
	curSelected := 0
	moved := 0
	
	//Service each request in order of SSTF
	for i := 0; i < initLength; i++ {
		//We first must find the shortest seek time
		for j := 0; j < len(requestArr); j++ {
			moved = int(math.Abs(float64(curCyl - requestArr[j])))
			if(moved < curShortest){
				curShortest = moved
				curSelected = j
			}
		}
		
		//After we find the shortest seek time, service it if it's a valid request
		if((requestArr[curSelected] > upperCyl) || (requestArr[curSelected] < lowerCyl)) {
			fmt.Printf("ERROR(15):Request out of bounds:  req (%d) > upper (%d) or < lower (%d)\n", requestArr[curSelected], upperCyl, lowerCyl)
			
			//If it isn't valid, we still need to remove it from the list
			requestArr = append(requestArr[:curSelected], requestArr[(curSelected+1):]...)
			continue
		}
		
		//Once we find the closest request, service it
		//No need to calculate distance traversed a second time
		fmt.Printf("Servicing %5d\n", requestArr[curSelected])
		
		//After we've serviced the request, remove the request from the slice and reset variables
		curCyl = requestArr[curSelected]
		requestArr = append(requestArr[:curSelected], requestArr[curSelected+1:]...)
		traverse = traverse + int(curShortest)
		curShortest = math.MaxInt32
		curSelected = 0
	}
	
	fmt.Printf("SSTF traversal count = %5d\n", traverse)
}

func scan(reqArr []int, initCyl int, lowerCyl int, upperCyl int) {
	//First we need to do a data dump before procession the scheduling
	fmt.Printf("Seek algorithm: SCAN\n")
	fmt.Printf("\tLower cylinder: %5d\n", lowerCyl)
	fmt.Printf("\tUpper cylinder: %5d\n", upperCyl)
	fmt.Printf("\tInit cylinder:  %5d\n", initCyl)
	fmt.Printf("\tCylinder requests:\n")
	
	for i := 0; i < len(reqArr); i++ {
		fmt.Printf("\t\tCylinder %5d\n", reqArr[i])
	}

	//Let's use a custom sorting function because Eustis still isn't updated...
	MergeSort(reqArr)
	
	length := len(reqArr)
	traverse := 0
	moved := 0
	curCyl := initCyl
	partition := len(reqArr)
	
	//Partition the array into 2 parts: larger than initCyl and smaller than initCyl
	for i := 0; i < len(reqArr); i++ {
		if (reqArr[i] >= initCyl){
			partition = i
			break
		}
	}
	
	//Scan upwards from the initial cylinder
	for i := partition; i < length; i++ {
		//make sure the request is valid
		if((reqArr[i] > upperCyl) || (reqArr[i] < lowerCyl)) {
			fmt.Printf("ERROR(15):Request out of bounds:  req (%d) > upper (%d) or < lower (%d)\n", reqArr[i], upperCyl, lowerCyl)
			continue
		}
		
		fmt.Printf("Servicing %5d\n", reqArr[i])
		moved = int(math.Abs(float64(curCyl - reqArr[i])))
		traverse = traverse + moved
		curCyl = reqArr[i]
	}
	
	//After we scan up to the highest request, move to the upper cylinder if there are requests left
	if partition > 0 {
		moved = int(math.Abs(float64(curCyl - upperCyl)))
		traverse = traverse + moved
	}
	curCyl = upperCyl
	
	//Now start scanning back down
	for i := (partition - 1); i >= 0; i-- {
		//make sure the request is valid
		if((reqArr[i] > upperCyl) || (reqArr[i] < lowerCyl)) {
			fmt.Printf("ERROR(15):Request out of bounds:  req (%d) > upper (%d) or < lower (%d)\n", reqArr[i], upperCyl, lowerCyl)
			continue
		}
		
		fmt.Printf("Servicing %5d\n", reqArr[i])
		moved = int(math.Abs(float64(curCyl - reqArr[i])))
		traverse = traverse + moved
		curCyl = reqArr[i]
	}
	
	fmt.Printf("SCAN traversal count = %5d\n", traverse)
}

func cscan(reqArr []int, initCyl int, lowerCyl int, upperCyl int) {
	//First we need to do a data dump before procession the scheduling
	fmt.Printf("Seek algorithm: C-SCAN\n")
	fmt.Printf("\tLower cylinder: %5d\n", lowerCyl)
	fmt.Printf("\tUpper cylinder: %5d\n", upperCyl)
	fmt.Printf("\tInit cylinder:  %5d\n", initCyl)
	fmt.Printf("\tCylinder requests:\n")
	
	for i := 0; i < len(reqArr); i++ {
		fmt.Printf("\t\tCylinder %5d\n", reqArr[i])
	}

	//Let's use a custom sorting function because Eustis still isn't updated...
	MergeSort(reqArr)
	
	length := len(reqArr)
	traverse := 0
	moved := 0
	curCyl := initCyl
	partition := len(reqArr)
	
	//Partition the array into 2 parts: larger than initCyl and smaller than initCyl
	for i := 0; i < len(reqArr); i++ {
		if (reqArr[i] >= initCyl){
			partition = i
			break
		}
	}
	
	//Scan upwards from the initial cylinder
	for i := partition; i < length; i++ {
		//make sure the request is valid
		if((reqArr[i] > upperCyl) || (reqArr[i] < lowerCyl)) {
			fmt.Printf("ERROR(15):Request out of bounds:  req (%d) > upper (%d) or < lower (%d)\n", reqArr[i], upperCyl, lowerCyl)
			continue
		}
		
		fmt.Printf("Servicing %5d\n", reqArr[i])
		
		moved = int(math.Abs(float64(curCyl - reqArr[i])))
		traverse = traverse + moved
		curCyl = reqArr[i]
	}
	
	//After we scan up to the highest request, move to the upper cylinder, then reset to the lower
	//Only if there are requests left
	if partition > 0 {
		moved = int(math.Abs(float64(curCyl - upperCyl)))
		traverse = traverse + moved
		traverse = traverse + upperCyl //Account for moving to the lower cylinder
	}
	
	curCyl = lowerCyl
	
	//Now start scanning back up from the bottom
	for i := 0; i < partition; i++ {
		//make sure the request is valid
		if((reqArr[i] > upperCyl) || (reqArr[i] < lowerCyl)) {
			fmt.Printf("ERROR(15):Request out of bounds:  req (%d) > upper (%d) or < lower (%d)\n", reqArr[i], upperCyl, lowerCyl)
			continue
		}
		
		fmt.Printf("Servicing %5d\n", reqArr[i])
		moved = int(math.Abs(float64(curCyl - reqArr[i])))
		traverse = traverse + moved
		curCyl = reqArr[i]
	}
	
	fmt.Printf("C-SCAN traversal count = %5d\n", traverse)
}

func look(reqArr []int, initCyl int, lowerCyl int, upperCyl int){
	//First we need to do a data dump before procession the scheduling
	fmt.Printf("Seek algorithm: LOOK\n")
	fmt.Printf("\tLower cylinder: %5d\n", lowerCyl)
	fmt.Printf("\tUpper cylinder: %5d\n", upperCyl)
	fmt.Printf("\tInit cylinder:  %5d\n", initCyl)
	fmt.Printf("\tCylinder requests:\n")
	
	for i := 0; i < len(reqArr); i++ {
		fmt.Printf("\t\tCylinder %5d\n", reqArr[i])
	}

	//Let's use a custom sorting function because Eustis still isn't updated...
	MergeSort(reqArr)
	
	length := len(reqArr)
	traverse := 0
	moved := 0
	curCyl := initCyl
	partition := len(reqArr)
	
	//Partition the array into 2 parts: larger than initCyl and smaller than initCyl
	for i := 0; i < len(reqArr); i++ {
		if (reqArr[i] >= initCyl){
			partition = i
			break
		}
	}
	
	//Scan upwards from the initial cylinder
	for i := partition; i < length; i++ {
		//make sure the request is valid
		if((reqArr[i] > upperCyl) || (reqArr[i] < lowerCyl)) {
			fmt.Printf("ERROR(15):Request out of bounds:  req (%d) > upper (%d) or < lower (%d)\n", reqArr[i], upperCyl, lowerCyl)
			continue
		}
		
		fmt.Printf("Servicing %5d\n", reqArr[i])
		moved = int(math.Abs(float64(curCyl - reqArr[i])))
		traverse = traverse + moved
		curCyl = reqArr[i]
	}
	
	//After we scan up to the highest request, scan down
	for i := (partition - 1); i >= 0; i-- {
		//make sure the request is valid
		if((reqArr[i] > upperCyl) || (reqArr[i] < lowerCyl)) {
			fmt.Printf("ERROR(15):Request out of bounds:  req (%d) > upper (%d) or < lower (%d)\n", reqArr[i], upperCyl, lowerCyl)
			continue
		}
		
		fmt.Printf("Servicing %5d\n", reqArr[i])
		moved = int(math.Abs(float64(curCyl - reqArr[i])))
		traverse = traverse + moved
		curCyl = reqArr[i]
	}
	
	fmt.Printf("LOOK traversal count = %5d\n", traverse)
}

func clook(reqArr []int, initCyl int, lowerCyl int, upperCyl int) {
	//First we need to do a data dump before procession the scheduling
	fmt.Printf("Seek algorithm: C-LOOK\n")
	fmt.Printf("\tLower cylinder: %5d\n", lowerCyl)
	fmt.Printf("\tUpper cylinder: %5d\n", upperCyl)
	fmt.Printf("\tInit cylinder:  %5d\n", initCyl)
	fmt.Printf("\tCylinder requests:\n")
	
	for i := 0; i < len(reqArr); i++ {
		fmt.Printf("\t\tCylinder %5d\n", reqArr[i])
	}

	//Let's use a custom sorting function because Eustis still isn't updated...
	MergeSort(reqArr)
	
	length := len(reqArr)
	traverse := 0
	moved := 0
	curCyl := initCyl
	partition := len(reqArr)
	
	//Partition the array into 2 parts: larger than initCyl and smaller than initCyl
	for i := 0; i < len(reqArr); i++ {
		if (reqArr[i] >= initCyl){
			partition = i
			break
		}
	}
	
	//Scan upwards from the initial cylinder
	for i := partition; i < length; i++ {
		//make sure the request is valid
		if((reqArr[i] > upperCyl) || (reqArr[i] < lowerCyl)) {
			fmt.Printf("ERROR(15):Request out of bounds:  req (%d) > upper (%d) or < lower (%d)\n", reqArr[i], upperCyl, lowerCyl)
			continue
		}
		
		fmt.Printf("Servicing %5d\n", reqArr[i])
		moved = int(math.Abs(float64(curCyl - reqArr[i])))
		traverse = traverse + moved
		curCyl = reqArr[i]
	}
	
	//After we scan up to the highest request, reset to the lower
	curCyl = reqArr[0]
	
	//Only move to the bottom to scan if there are more requests left
	if partition > 0{
		moved = int(math.Abs(float64(curCyl - reqArr[(len(reqArr)-1)])))
		traverse = traverse + moved
	}
	
	
	//Now start scanning back up
	for i := 0; i < partition; i++ {
		//make sure the request is valid
		if((reqArr[i] > upperCyl) || (reqArr[i] < lowerCyl)) {
			fmt.Printf("ERROR(15):Request out of bounds:  req (%d) > upper (%d) or < lower (%d)\n", reqArr[i], upperCyl, lowerCyl)
			continue
		}
		
		fmt.Printf("Servicing %5d\n", reqArr[i])
		moved = int(math.Abs(float64(curCyl - reqArr[i])))
		traverse = traverse + moved
		curCyl = reqArr[i]
	}
	
	fmt.Printf("C-LOOK traversal count = %5d\n", traverse)
}

func main(){
	//Grab the input file name from the command land
	//Open up the input file 
	inp, err := os.Open(os.Args[1])
	if err != nil{
		log.Fatal(err)
	}
	
	defer inp.Close()
	
	//Declare variables to prevent scoping issues
	lowerCyl := 0
	upperCyl := 0
	initCyl := 0
	schedAlg := 0
	requestedCyl := 0
	var requestArr []int
	
	//Parse the file input
	//For this approach, we're going to read in all of the file before we begin processing
	scanner := bufio.NewScanner(inp)
	for scanner.Scan() {
		curLine := scanner.Text()
		
		//After we scan the line of text, we need to parse it
		parsedLine := strings.Fields(curLine)
		
		//Time to try to deal with each word in the line
		for i := 0; i < len(parsedLine); i++{
			//If we find "cylreq" we need to append a new cylinder request to the slice
			if parsedLine[i] == "cylreq" {
				i += 1
				test := parsedLine[i]
				requestedCyl, err = strconv.Atoi(test)
				if err != nil { //error handling
					fmt.Println("Input error, cylreq must be an integer. Ending")
					os.Exit(3)
				}
				
				requestArr = append(requestArr, requestedCyl)
				
				break
			}
			
			//If we find "initCYL" we need to define the initial cylinder for the simulation
			if parsedLine[i] == "initCYL" {
				i += 1
				test := parsedLine[i]
				initCyl, err = strconv.Atoi(test)
				if err != nil { //error handling
					fmt.Println("Input error, initCyl must be an integer. Ending")
					os.Exit(3)
				}
				
				break
			}
			
			//If we find "use" we need see which algorithm we're using
			if parsedLine[i] == "use" {
				i += 1
				
				//Time to check which scheduling algorithm was found
				if parsedLine[i] == "fcfs" {
					schedAlg = 1
				}
				
				if parsedLine[i] == "sstf" {
					schedAlg = 2
				}
				
				if parsedLine[i] == "scan" {
					schedAlg = 3
				}
				
				if parsedLine[i] == "c-scan" {
					schedAlg = 4
				}
				
				if parsedLine[i] == "look" {
					schedAlg = 5
				}
				
				if parsedLine[i] == "c-look" {
					schedAlg = 6
				}

				break
			}
			
			//If we find "lowerCYL" we need to define the lowest cylinder limit
			if parsedLine[i] == "lowerCYL" {
				i += 1
				test := parsedLine[i]
				lowerCyl, err = strconv.Atoi(test)
				if err != nil { //error handling
					fmt.Println("Input error, lowerCYL must be an integer. Ending")
					os.Exit(3)
				}
				
				break
			}
			
			//If we find "upperCYL" we need to define the highest cylinder limit
			if parsedLine[i] == "upperCYL" {
				i += 1
				test := parsedLine[i]
				upperCyl, err = strconv.Atoi(test)
				if err != nil { //error handling
					fmt.Println("Input error, upperCYL must be an integer. Ending")
					os.Exit(3)
				}
				
				break
			}
			
			//If we find "end" then we're at the last line of the input
			if parsedLine[i] == "end" {
				break
			}
			
			break
		}
	}
	
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	
	//We need to check if the initial cylinder, upper cylinder, and lower cylinder positions are valid
	if initCyl > upperCyl {
		fmt.Printf("ABORT(11):initial (%d) > upper (%d)\n", initCyl, upperCyl)
		os.Exit(0)
	}
	
	if initCyl < lowerCyl {
		fmt.Printf("ABORT(12):initial (%d) < lower (%d)\n", initCyl, lowerCyl)
		os.Exit(0)
	}
	
	if lowerCyl > upperCyl {
		fmt.Printf("ABORT(13):upper (%d) < lower (%d)\n", upperCyl, lowerCyl)
		os.Exit(0)
	}
	
	//Once we've finished parsing the file input, we must begin scheduling
	//First one is fcfs
	if schedAlg == 1 {
		fcfs(requestArr, initCyl, lowerCyl, upperCyl)
	}
	
	//Second is sstf
	if schedAlg == 2 {
		sstf(requestArr, initCyl, lowerCyl, upperCyl)
	}
	
	//third is scan
	if schedAlg == 3 {
		scan(requestArr, initCyl, lowerCyl, upperCyl)
	}
	
	//fourth is c-scan
	if schedAlg == 4 {
		cscan(requestArr, initCyl, lowerCyl, upperCyl)
	}
	
	//fifth is look
	if schedAlg == 5{
		look(requestArr, initCyl, lowerCyl, upperCyl)
	}
	
	//Last is c-look
	if schedAlg == 6 {
		clook(requestArr, initCyl, lowerCyl, upperCyl)
	}
}