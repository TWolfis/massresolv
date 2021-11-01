//return a list of all hosts in an ip range with hostname
package main

import (
        "fmt"
        "net"
	"os"
	"encoding/binary"
	"time"
)

func main(){

	if len(os.Args) != 2{
		fmt.Println("Provide network/bits as input\n")
		os.Exit(1)
	}


	//parse ip to *ipnet struct
	_, ipv4Net,err := net.ParseCIDR(os.Args[1])
	if err != nil{
		fmt.Println(err)
		os.Exit(-1)
	}

	// convert IPNet struct mask and address to uint32
	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	start := binary.BigEndian.Uint32(ipv4Net.IP)
	finish := (start & mask) | (mask ^ 0xffffffff)
	
	startTime:= time.Now()
	//create channel
	ch := make(chan string)

	for i := start; i <= finish; i++{
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip,i)
		go resolv(ip.String(),ch)
	}

	//wait for output from channel 
	for i := start; i <= finish; i++{
		fmt.Printf(<-ch)
	}

	fmt.Printf("%.2fs elapsed\n",time.Since(startTime).Seconds())



}

//resolve ip address from range of hosts
func resolv(ip string, ch chan <- string){
	start := time.Now()
	names, err := net.LookupAddr(ip)
	
	//return empty output (so itterator does not get stuck)
	if err !=nil || len(names) == 0{
		ch <- fmt.Sprint()
		return
	}

 	stop:=time.Since(start).Seconds()
	ch <- fmt.Sprintf("%s resolved to %s in %.2fs\n",ip,names,stop)
}



