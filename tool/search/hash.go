//code for docHash

package search

//parses the data into blocks of 4 bits, while the maximum considered distance
//between values would be 16 bytes. Returns a bool array that with a length of 4096
func DocHash_4_16_4096(data []byte) []bool {
	sum:= make([]bool, 4096);
	pre:= make([]byte, 2*len(data));
	for i:=0; i<len(data); i++ {
		/*
		* split current byte into 2 peices
		* original:|<-first 4 bits->|<-last 4 bits->|
		*/
		// pre[i*2] = 0_|<-first 4 bits->|
		pre[i*2] = data[i]>>4;
		// pre[i*2+1] = 1_|<-last 4 bits->|
		pre[i*2+1] = data[i]%(1<<4) | (1<<4);
	}
	for i:=0; i<len(pre); i++ {
		for j:=0; j<32 && j<len(pre); j++ {
			dis:= func(x uint) int {
				if x < 2 {
					return 0;
				}
				if x < 4 {
					return 1;
				}
				if x < 8 {
					return 0b10;
				}
				if x < 16 {
					return 0b11;
				}
				return -1; //ERROR!
			}((uint)(j/2));
			// |<-2 bits distance->|<-5 bits data 1->|<-5 bits data 2->|
			id:= (dis<<10) | ((int)(pre[i])<<5) | (int)(pre[i+1+j]);
			sum[id] = true;
		}
	}
	return sum;
}