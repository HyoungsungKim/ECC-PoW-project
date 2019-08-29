#include "LDPC.h"
#include <string>
#include <iostream>
#include <ctime>
using namespace std;

int main(int argc, char* argv[]) {
	clock_t start, end;
	double elapseTime;

	start = clock();

	char phv1[32] = "0000000000000000000000000000000";
	const unsigned char phv[32] = "0000000000000000000000000000000";

	unsigned int nonce = 0;
	LDPC *ptr = new LDPC;

	ptr->set_difficulty(16,3,4);				//2 => n = 64, wc = 3, wr = 6, 	
	if (!ptr->initialization())
	{
		printf("error for calling the initialization function");
		return 0;
	}

	//Generate parity check matrix with previous hash vector.
	ptr->generate_seed(phv1);
	ptr->generate_H();
	ptr->generate_Q();
	ptr->generate_hv(phv);

	for(int i = 0; i < 100000; ++i) {	
		ptr->decoding();		
	}

	end = clock();
	elapseTime = (double)(end -start)/CLOCKS_PER_SEC;
	cout << elapseTime << "s";
}


