﻿// ECCPOW-LDPC.cpp: 콘솔 응용 프로그램의 진입점을 정의합니다.
//


#include "stdafx.h"
#include "LDPC.h"
#include <string>
#include <iostream>
using namespace std;

int main(int argc, char* argv[])
{
	//argv[0] : test.exe
	//argv[1] : previous block hash		-> 32 Hexa digit
	//argv[2] : current block header	-> 32 Hexa digit

	//phv : previous block hash
	//current_block_header : current block header
	string phv;
	string current_block_header_hash;
	if (argc == 1) {
		//const for test
		//원래 27개 였음
		phv = "00000000000000000000000000000000";
		current_block_header_hash = "00000000000000000000000000000000";
	}
	else {
		phv = argv[1];
		current_block_header_hash = argv[2];
	}

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

	unsigned int nonce = 0;
	LDPC *ptr = new LDPC;

	ptr->set_difficulty(24,3,6);				//2 => n = 64, wc = 3, wr = 6, 	
	if (!ptr->initialization())
	{
		printf("error for calling the initialization function");
		return 0;
	}
	
	char charPhv[32];
	copy(phv.begin(), phv.end(), charPhv);

	ptr->generate_seed(charPhv);
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
		bool flag = ptr->decision();
		if (!flag) // If a hash vector is a codeword itself, we dont need to run the decoding function.
		{
			ptr->decoding();
			flag = ptr->decision();
		}
		if (flag)
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

