package Assignment0203IBC

import (
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type Block struct {
	//Hash
	//Data

	Transaction string
	Reward string
	PrevPointer *Block
	Hash [32]byte
	PrevHash [32]byte
}

type Person struct{

	Name string
	Wallet float64
	Port string
	MinePort string
}

type Transaction struct{

	Sender int
	Amount float64
	Receiver int
	Miner int
	Transfee float64

}


func  DeriveHash(Transaction string)[32]byte {
	return sha256.Sum256([]byte(Transaction))
}


func InsertBlock(Transaction,Reward string, chainHead *Block) *Block {
	if chainHead == nil {
		return &Block{Transaction,Reward, nil,DeriveHash(Transaction+"+"+Reward),[32]byte{}}
	}
	return &Block{Transaction,Reward, chainHead,DeriveHash(Transaction+"+"+Reward),DeriveHash(chainHead.Transaction+"+"+chainHead.Reward)}

}

func ListBlocks(chainHead *Block) {
	for p := chainHead; p != nil; p = p.PrevPointer {
		fmt.Printf("Transaction: %s, Hash:%x, PrevHash:%x\n",p.Transaction+" & "+p.Reward,p.Hash,p.PrevHash)

	}
}




/*
func ChangeBlock(oldTrans string, newTrans string, chainHead *Block) {
	for p := chainHead; p != nil; p = p.next {
		if p.Transaction == oldTrans{
			p.Transaction=newTrans
			p.Hash=DeriveHash(newTrans)

		}
	}
}
*/
func VerifyChain(chainHead *Block) {
	for p := chainHead; p != nil; p = p.PrevPointer {
		if p.PrevPointer !=nil{
			if p.PrevHash != p.PrevPointer.Hash{
				fmt.Println("Chain is Invalid!")
				return
			}
		}
	}
	fmt.Println("Chain is Valid!")
}



func handleConnection(c net.Conn) {
	log.Println("A client has connected", c.RemoteAddr())
	c.Write([]byte("Hello world"))
}

func ExecuteTransaction(Sender *Person,SendingAmount float64,Receiver,Miner *Person,TransationFee,MinerReward float64,chainHead *Block)*Block{
	if Sender.Wallet>=SendingAmount{
		if TransationFee<=SendingAmount {
			Sender.Wallet -= SendingAmount
			Receiver.Wallet += SendingAmount - TransationFee
			Miner.Wallet += MinerReward+TransationFee
			return InsertBlock(Sender.Name+"-"+fmt.Sprintf("%f", SendingAmount)+" ==> "+fmt.Sprintf("%f", SendingAmount- TransationFee)+"->"+Receiver.Name, fmt.Sprintf("%f", MinerReward+TransationFee)+"->"+Miner.Name, chainHead)

		}else{
			fmt.Println("Invalid Reward Amount!")
		}
	}else{
		fmt.Println(Sender.Name,"does not have enough BC!")
	}
	return chainHead
}


func ReceiveUpdatedBlockChain(ln net.Listener,chainHead *Block,people []Person) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		dec := gob.NewDecoder(conn)
		err = dec.Decode(&chainHead)
		if err != nil {
			log.Println(err)
		}
		err = dec.Decode(&people)
		if err != nil {
			log.Println(err)
		}

		fmt.Println("Latest Blockchain and Ledger!")
		ListBlocks(chainHead)
		fmt.Println(people)
		fmt.Println("1. Do you want to make a transaction?\n2. Do you want to exit?")

	}
}
func BroadCastBlockChainAndLedger(chainHead *Block,people []Person,ind int){
	for i:=0;i<len(people);i++ {
		if i != ind {
			fmt.Println("Sending BlockChain & Ledger to " + people[i].Name)
			conn, err := net.Dial("tcp", "localhost:"+people[i].Port)
			if err != nil {
				//handle error
				fmt.Println(err)
			}
			gobEncoder := gob.NewEncoder(conn)
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
	}
}

func MinerHandler(ln net.Listener,chainHead *Block,people []Person,minerPort string, ind int) {
	ln, err := net.Listen("tcp", ":"+minerPort)
	if err != nil {
		log.Fatal(err)
	}
	MinerReward:=100.0
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		var transaction Transaction
		gobEncoder := gob.NewDecoder(conn)
		err = gobEncoder.Decode(&transaction)
		fmt.Println(transaction)
		chainHead = ExecuteTransaction(&people[transaction.Sender],transaction.Amount,&people[transaction.Receiver],&people[transaction.Miner],transaction.Transfee,MinerReward,chainHead)
		MinerReward-=10
		fmt.Println("Latest Blockchain and Ledger!")
		ListBlocks(chainHead)
		fmt.Println(people)
		fmt.Println("BroadCasting!")
		BroadCastBlockChainAndLedger(chainHead,people,ind)
		fmt.Println("1. Do you want to make a transaction?\n2. Do you want to exit?")
	}
}
