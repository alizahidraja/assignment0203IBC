package main

import (
	a2 "Assignment0203IBC"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
)

func main() {
	/*
		fmt.Println(reflect.TypeOf(tst))
		tst1 := &a2.Person{"",0.0,""}

		fmt.Println(reflect.TypeOf(tst1))
		fmt.Println(reflect.TypeOf(tst1)==reflect.TypeOf(&a2.Person{"",0,""}))
	*/

	fmt.Println("Welcone to BlockCoin (BC), The developer of this Currency is Ali")
	MinerReward:=100.0
	if len(os.Args) == 2{//Satoshi
		//max users=os.args[1]
		Users,_:=strconv.Atoi(os.Args[1])
		var chainHead *a2.Block
		//var max_nodes int
		//_, _ = fmt.Scan(&age)

		//reader := bufio.newReader(os.Stdin)
		//var name string
		//fmt.Println("What is your name?")
		//name, _ := reader.readString("\n")

		//max_nodes = 6
		//_, _ = fmt.Scan(&max_nodes)
		//var nodes int = 0

		people:=make([]a2.Person,0)
		people = append(people,a2.Person{"Ali",100,"6000","5000"} )
		//people = append(people,a2.Person{"Dani",0,6001} )
		//people = append(people,a2.Person{"Saad",0,6001} )
		//people = append(people,a2.Person{"Shah",0,6001} )
		//people = append(people,a2.Person{"Joan",0,6001} )
		//people = append(people,a2.Person{"Raja",0,6001} )

		fmt.Println(people)

		fmt.Println("This is Ali, making the genesis Block")

		chainHead = a2.InsertBlock("GenesisBlock", "100->Ali",nil)

		//chainHead = Executea2.Transaction(&people[0],2,&people[1],&people[2],0.2,chainHead)
		//chainHead = Executea2.Transaction(&people[1],2,&people[3],&people[4],0,chainHead)
		//chainHead = Executea2.Transaction(&people[0],2,&people[1],&people[2],0.2,chainHead)
		//chainHead = Executea2.Transaction(&people[1],2,&people[3],&people[4],0,chainHead)
		//chainHead = Executea2.Transaction(&people[0],2,&people[4],&people[5],0.2,chainHead)
		//chainHead = Executea2.Transaction(&people[0],2,&people[5],&people[4],0.2,chainHead)

		fmt.Println(people)
		a2.ListBlocks(chainHead)
		//ChangeBlock("AliceToBob", "AliceToTrudy", chainHead)
		a2.VerifyChain(chainHead)
		//		fmt.Println(len(os.Args), os.Args)



		ln, err := net.Listen("tcp", ":6000")
		if err != nil {
			log.Fatal(err)
		}
		for {//Waiting for All Users
			fmt.Println("Waiting for Users")
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			//gobEncoder := gob.NewEncoder(conn)
			//err = gobEncoder.Encode(chainHead)
			var person a2.Person
			gobDecoder:=gob.NewDecoder(conn)
			err=gobDecoder.Decode(&person)
			fmt.Println("New a2.Person "+person.Name+" added!")
			people=append(people,person)
			chainHead = a2.ExecuteTransaction(&people[len(people)-1],0,&people[len(people)-1],&people[0],0,MinerReward,chainHead)
			if err != nil {
				log.Println(err)
			}
			Users=Users-1
			if Users==0{
				fmt.Println("Max users have entered the network, starting propagation!")
				break
			}
		}
		fmt.Println(people)

		for i:=1;i<len(people);i++ {

			fmt.Println("Sending BlockChain & Ledger to " + people[i].Name)
			conn, err := net.Dial("tcp", "localhost:"+people[i].Port)
			if err != nil {
				//handle error
				fmt.Println(err)
			}
			gobEncoder := gob.NewEncoder(conn)
			err=gobEncoder.Encode(i)
			err = gobEncoder.Encode(chainHead)
			//gob.RegisterName()
			if err != nil {
				//handle error
				fmt.Println(err)
			}
			err = gobEncoder.Encode(people)
			if err != nil {
				//handle error
				fmt.Println(err)
			}

		}
		fmt.Println("Latest Blockchain!")
		a2.ListBlocks(chainHead)

		go a2.ReceiveUpdatedBlockChain(ln,chainHead,people)
		go a2.MinerHandler(ln,chainHead,people,"5000",0)
		for{
			fmt.Println("1. Do you want to make a Transaction?\n2. Do you want to exit?")
			var i,j int
			var amount,tranfee float64
			_, err := fmt.Scanln(&i)
			if err!=nil{
				fmt.Println(err)
			}
			if i == 1{//a2.Transaction
				fmt.Println("Who do you want to send BC to?")
				for i=0;i<len(people);i++{
					if i!=0 {
						fmt.Println(i,people[i].Name)
					}
				}
				_, err := fmt.Scanln( &j)
				if err!=nil{
					fmt.Println(err)
				}

				fmt.Println("How much do you want to send to",people[j].Name,"?")
				_, err = fmt.Scanln( &amount)
				if err!=nil{
					fmt.Println(err)
				}

				fmt.Println("How much do you wanna give as Transaction fees ?")
				_, err = fmt.Scanln( &tranfee)
				if err!=nil{
					fmt.Println(err)
				}

				miner:=rand.Intn(len(people))
				fmt.Println(miner,people[miner],"is selected as miner")

				fmt.Println("Sending Transaction to Miner, i.e.",people[miner].Name)
				conn, err := net.Dial("tcp", "localhost:"+people[miner].MinePort)
				if err != nil {
					//handle error
					fmt.Println(err)
				}
				gobEncoder := gob.NewEncoder(conn)
				err = gobEncoder.Encode(a2.Transaction{0,amount,j,miner,tranfee})

				/*chainHead = Executea2.Transaction(&people[0],amount,&people[j],&people[miner],tranfee,chainHead)
				fmt.Println("Latest Blockchain and Ledger!")
				a2.ListBlocks(chainHead)
				fmt.Println(people)
				fmt.Println("BroadCasting!")
				broadCastBlockChainAndLedger(chainHead,people,0)
				*/
			}else if i ==2{
				break
			}
		}

	}else {//Peers
		//arg[1]=Name, arg[2]=Port, arg[3]=MinePort
		fmt.Println("This is a new member of the block, named ",os.Args[1])
		conn, err := net.Dial("tcp", "localhost:6000")
		if err != nil {
			//handle error
		}
		person:=&a2.Person{os.Args[1],0,os.Args[2],os.Args[3]}
		enc :=gob.NewEncoder(conn)
		err=enc.Encode(person)




		ln, err := net.Listen("tcp", ":"+os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		conn, err = ln.Accept()
		if err != nil {
			log.Println(err)
		}
		var ind int
		var chainHead *a2.Block
		dec := gob.NewDecoder(conn)
		err=dec.Decode(&ind)

		err = dec.Decode(&chainHead)
		if err != nil {
			log.Println(err)
		}
		var people []a2.Person
		err = dec.Decode(&people)
		if err != nil {
			log.Println(err)
		}

		a2.ListBlocks(chainHead)
		fmt.Println(people)

		go a2.ReceiveUpdatedBlockChain(ln,chainHead,people)
		go a2.MinerHandler(ln,chainHead,people,os.Args[3],ind)
		for{
			fmt.Println("1. Do you want to make a Transaction?\n2. Do you want to exit?")
			var i,j int
			var amount,tranfee float64
			_, err := fmt.Scanln(&i)
			if err!=nil{
				fmt.Println(err)
			}
			if i == 1{//a2.Transaction
				fmt.Println("Who do you want to send BC to?")
				for i=0;i<len(people);i++{
					if i!=ind {
						fmt.Println(i,people[i].Name)
					}
				}
				_, err := fmt.Scanln( &j)
				if err!=nil{
					fmt.Println(err)
				}

				fmt.Println("How much do you want to send to",people[j].Name,"?")
				_, err = fmt.Scanln( &amount)
				if err!=nil{
					fmt.Println(err)
				}

				fmt.Println("How much do you wanna give as Transaction fees ?")
				_, err = fmt.Scanln( &tranfee)
				if err!=nil{
					fmt.Println(err)
				}

				miner:=rand.Intn(len(people))
				fmt.Println(miner,people[miner],"is selected as miner")
				fmt.Println("Sending a2.Transaction to Miner, i.e.",people[miner])
				conn, err := net.Dial("tcp", "localhost:"+people[miner].MinePort)
				if err != nil {
					//handle error
					fmt.Println(err)
				}
				gobEncoder := gob.NewEncoder(conn)
				err = gobEncoder.Encode(a2.Transaction{ind,amount,j,miner,tranfee})

				/*chainHead = Executea2.Transaction(&people[ind],amount,&people[j],&people[miner],tranfee,chainHead)
				fmt.Println("Latest Blockchain and Ledger!")
				a2.ListBlocks(chainHead)
				fmt.Println(people)
				fmt.Println("BroadCasting!")
				broadCastBlockChainAndLedger(chainHead,people,ind)
				*/

			}else if i ==2{
				break
			}
		}

	}

}

