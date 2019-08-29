#include "LDPC.h"
#include <string>
#include <iostream>
using namespace std;

int main(int argc, char* argv[])
{
	string phv;
	string current_block_header_hash;
	if (argc == 1) {
		phv = "00000000000000000000000000000000";
		current_block_header_hash = "00000000000000000000000000000000";
	}
	else {
		phv = argv[1];
		current_block_header_hash = argv[2];
	}
	//Resize Warning! - Only for test
	phv.resize(32, '0');
	//current_block_header does not need to resize. Resize can cause hash collusion. So comment out it
	//current_block_header_hash.resize(32, '0');

	/*
	if (phv.size() < 32) {
		while (phv.size() < 32) {
			phv += '0';
		}
	}

	if (current_block_header_hash.size() < 32) {
		while (current_block_header_hash.size() < 32) {
			current_block_header_hash += '0';
		}
	}
	*/

	unsigned int nonce = 0;
	LDPC *ptr = new LDPC;


	ptr->set_difficulty(24,3,6);				//2 => n = 64, wc = 3, wr = 6, 	
	if (!ptr->initialization())
	{
		printf("error for calling the initialization function");
		return 0;
	}

	//Generate parity check matrix with previous hash vector.
	ptr->generate_seed(phv);
	ptr->generate_H();
	ptr->generate_Q();
	
	ptr->print_H("H2.txt");
	ptr->print_Q(NULL, 1);
	ptr->print_Q(NULL, 2);


	while (1)
	{
		string current_block_header_with_nonce;
		current_block_header_with_nonce.assign(current_block_header_hash);
		current_block_header_with_nonce += to_string(nonce);

		ptr->generate_hv((unsigned char*)current_block_header_with_nonce.c_str());
		//codeword 여부 판단 if codeword, flag is 1.
		bool flag = ptr->decision();		
		if (!flag) //If flag is 0(false), do decoding again.
			// If a hash vector is a codeword itself, we dont need to run the decoding function.
		{
			ptr->decoding();
			flag = ptr->decision();
		}
		if (flag)	//If flag is 1(true), Decoding is completed.
		{
			printf("codeword is founded with nonce = %d\n", nonce);
			FILE* fp;
			const char name[] = "nonceAndBlockHash.txt";
			if (name) {
				fopen_s(&fp, name, "w");
			}
			else {
				fp = stdout;
			}
		
			fprintf(fp,"%d\n", nonce);
			if (name)
				fclose(fp);

			break;
		
		}		
		nonce++;		
	}


	ptr->print_word(NULL, 1);
	ptr->print_word(NULL, 2);
	delete ptr;
	return 0;
}

