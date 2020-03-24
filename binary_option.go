package main

import (
	"fmt"
	"github.com/spf13/cast"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"
)


func main(){

	file, err := os.Create("bo_file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()


	file_1, err := os.Create("special_file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file_1.Close()


	round := 50000
	output_num_list := []int{}

	odd_rounds := 0
	even_rounds := 0
	under_rounds := 0
	over_rounds := 0

	max_odd_rounds := 0
	max_even_rounds := 0
	max_under_rounds := 0
	max_over_rounds := 0


	for i:=0;i<round;i++{

		probability_map := map[string]int{}
		bet_map,probability_total := gameSetUp(probability_map)

		//fmt.Println(bet_map)
		//fmt.Println(probability_map)

		output_num,output_region := randomOutputNum(probability_map,probability_total)
		banker_income := 0.0
		single_round_bet := incomeCal(bet_map,output_region,&banker_income)


		switch output_region{
		case "under_odd":
			under_rounds += 1
			odd_rounds += 1
			over_rounds = 0
			even_rounds = 0
			file_1.WriteString("-100 1"+"\n")

		case "under_even":
			under_rounds += 1
			even_rounds += 1
			odd_rounds = 0
			over_rounds = 0

			file_1.WriteString("-100 -1"+"\n")
		case "over_odd":
			under_rounds = 0
			odd_rounds += 1
			over_rounds += 1
			even_rounds = 0
			file_1.WriteString("100 1"+"\n")
		case "over_even":
			under_rounds = 0
			odd_rounds = 0
			over_rounds += 1
			even_rounds += 1
			file_1.WriteString("100 -1"+"\n")
		}

		if under_rounds > max_under_rounds{
			max_under_rounds = under_rounds
		}
		if odd_rounds >max_odd_rounds{
			max_odd_rounds = odd_rounds
		}
		if over_rounds > max_over_rounds{
			max_over_rounds = over_rounds
		}
		if even_rounds > max_even_rounds{
			max_even_rounds = even_rounds
		}

		//fmt.Println(output_num,output_region,single_round_bet)
		file.WriteString(cast.ToString(output_num) + " "+ cast.ToString(output_region) + " "+cast.ToString(single_round_bet) + " " +cast.ToString(banker_income) +"\n" )
		output_num_list = append(output_num_list,output_num)
	}

	sort.Ints(output_num_list)

	for i:= 0;i<len(output_num_list);i++{
		file.WriteString(cast.ToString(output_num_list[i])+"\n")

	}

	fmt.Println("max_over_rounds= ",max_over_rounds )
	fmt.Println("max_under_rounds=" ,max_under_rounds)
	fmt.Println("max_odd_rounds= ",max_odd_rounds)
	fmt.Println("max_even_rounds= ",max_even_rounds)

	//計算連續不開的機率為多少

}

func incomeCal(bet_map []kv,output_region string , banker_income *float64)float64{

	single_round_bet := 0.0

	for i:=0;i<len(bet_map);i++{

		single_round_bet += bet_map[i].Value

		if  bet_map[i].Key == output_region{
			*banker_income -= 1.82* bet_map[i].Value
		}

	}

	single_round_bet /= 2.0
	*banker_income += single_round_bet
	return single_round_bet
}


func betSetUP()(odd_total_bet,even_total_bet,over_total_bet,under_total_bet  float64){

	odd_player := []Player{}
	even_player := []Player{}
	over_player := []Player{}
	under_player := []Player{}
	player_mum := 30

	rand.Seed(time.Now().UnixNano())
	for i:= 0;i<player_mum;i++{
		bet_player := cast.ToFloat64(100*(rand.Intn(9)+1))
		pick_game := rand.Intn(4)
		switch pick_game{
		case 0:
			odd_player = append(odd_player,Player{bet_player,0})
			odd_total_bet += bet_player
		case 1:
			even_player = append(even_player,Player{bet_player,0})
			even_total_bet += bet_player
		case 2:
			over_player = append(over_player,Player{bet_player,0})
			over_total_bet += bet_player
		case 3:
			under_player = append(under_player,Player{bet_player,0})
			under_total_bet += bet_player
		}
	}
	//fmt.Println(odd_total_bet )
	//fmt.Println(even_total_bet )
	//fmt.Println(over_total_bet )
	//fmt.Println(under_total_bet )
	return
}

func gameSetUp(probability_map map[string]int)(bet_map []kv,probability_total int){


	odd_total_bet,even_total_bet,over_total_bet,under_total_bet := betSetUP()
	remain_probability := 1000

	//var bet_map []kv

	bet_map= append(bet_map, kv{"over_odd",over_total_bet + odd_total_bet})
	bet_map= append(bet_map, kv{"over_even",over_total_bet + even_total_bet})
	bet_map= append(bet_map, kv{"under_odd",under_total_bet + odd_total_bet})
	bet_map= append(bet_map, kv{"under_even",under_total_bet + even_total_bet})

	//bet_map= append(bet_map, kv{"over_odd",5000})
	//bet_map= append(bet_map, kv{"over_even",1000})
	//bet_map= append(bet_map, kv{"under_odd",20})
	//bet_map= append(bet_map, kv{"under_even",10})
/*
	bet_map=append(bet_map,kv{"over_odd",7400})
	bet_map=append(bet_map,kv{"over_even",4100})
	bet_map=append(bet_map,kv{"under_odd",10100})
	bet_map=append(bet_map,kv{"under_even",6800})
*/

	sort.Slice(bet_map, func(i, j int) bool {
		return bet_map[i].Value > bet_map[j].Value  // 降序
		// return ss[i].Value > ss[j].Value  // 升序
	})

	for i:=0;i<len(bet_map);i++{
		if isProbabilityValid(bet_map[i:],probability_map,&remain_probability){
			break
		}
	}
	probability_total = 1000 - remain_probability
	return
}


type kv struct {
	Key   string
	Value float64
}

func randomOutputNum(probability_map map[string]int, probability_total int)(int,string){

	rand.Seed(time.Now().UnixNano())
	probability_distribution := make([]string,probability_total)
	count  := 0
	for i,_ := range probability_map{
		for j:=0;j<probability_map[i];j++{
			probability_distribution[count] = i
			count += 1
		}
	}

	tmp := rand.Intn(probability_total)
	selected_region := probability_distribution[tmp]
	//fmt.Printf("selected number = %v ",selected_region)

	region_num_map := map[string][]int{
		"over_odd":{5,7,9},
		"over_even":{6,8},
		"under_odd":{1,3},
		"under_even":{0,2,4},
	}

	selected_num := region_num_map[selected_region][rand.Intn(len(region_num_map[selected_region]))]
	return selected_num,selected_region
}

func isProbabilityValid(bet_map []kv,probability_map map[string]int,remain_probability *int)bool{

	sum_bet := 0.0
	for i:=0;i<len(bet_map);i++{
		sum_bet += bet_map[i].Value
	}
	//avg_bet := sum_bet / cast.ToFloat64(len(bet_map))


	probability_tmp := setProbability(len(bet_map))

	if (bet_map[0].Value / sum_bet) > (probability_tmp[0]){
		probability_map[bet_map[0].Key] = cast.ToInt(math.Round(probability_tmp[0] * cast.ToFloat64(*remain_probability)))
		*remain_probability -= probability_map[bet_map[0].Key]

	}else{
		for i:=0;i<len(bet_map);i++{
			probability_map[bet_map[i].Key] = cast.ToInt(math.Round((bet_map[i].Value / sum_bet)*(cast.ToFloat64(*remain_probability))))
		}
		for i:=0;i<len(bet_map);i++{
			*remain_probability -= probability_map[bet_map[i].Key]
		}
		return true
	}
	return false
}


type Player struct{
	bet_money float64
	//bet_game int
	income float64
}

func setProbability(blocker_num int)([]float64){

	dev_container := []float64{}
	probability_arr := []float64{}

	for i:= 0;i<blocker_num+1;i++{
		dev_container = append(dev_container,-2.0 + 4.0/(cast.ToFloat64(blocker_num))* cast.ToFloat64(i))
	}

	mean := 0.0
	dev := 1.0

	for i:=1;i<len(dev_container);i++{

		interval := 0.001
		X := dev_container[i-1]+interval*(dev_container[i] - dev_container[i-1])
		sum := 0.0
		last_p := norDistribution(mean,dev,dev_container[i-1])
		now_p:= 0.0

		for j:=0;j<1000;j++{
			now_p = norDistribution(mean,dev,X)
			sum +=(now_p +last_p )/2 *interval*(dev_container[i]- dev_container[i-1])
			X += interval*(dev_container[i] - dev_container[i-1])
			last_p = now_p
		}
		probability_arr = append(probability_arr,sum)
	}

	sort.Float64s(probability_arr)
	for i, j := 0, len(probability_arr)-1; i < j; i, j = i+1, j-1 {
		probability_arr[i], probability_arr[j] = probability_arr[j], probability_arr[i]
	}
	return probability_arr
}

func norDistribution(mean,dev,X float64)float64{

	p := 1/(math.Sqrt(2*math.Pi)*dev)*math.Exp(-0.5*math.Pow((X-mean)/dev,2))
	return p
}