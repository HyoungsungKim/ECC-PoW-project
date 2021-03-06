// ECCPOW-LDPC.cpp: 콘솔 응용 프로그램의 진입점을 정의합니다.
//


#include "stdafx.h"
#include "LDPC.h"
#include <string>
#include <time.h>

using namespace std;

#define diff_lpdc_generate 0
int main()
{
	ldpc_level_parameter tmp;
#if diff_lpdc_generate
	FILE *fp; fopen_s(&fp, (const char*)"diff_ldpc.txt", "r");
	FILE *fp2; fopen_s(&fp2, (const char*)"diff_ldpc_fin.txt", "w");
	while (fp)
	{
		//fscanf_s(fp,)
		fscanf_s(fp, "%d,%d,%d,%d,%d,%d,%d,%lf,%le\n", &tmp.level, &tmp.n, &tmp.wc, &tmp.wr, &tmp.from, &tmp.to, &tmp.type, &tmp.clock, &tmp.prob);
		printf("%e\n", tmp.prob);
		fprintf_s(fp2, "{%d,\t%d,\t%d,\t%d,\t%d,\t%d,\t%d,\t%f,\t\t%e},\n", tmp.level, tmp.n, tmp.wc, tmp.wr, tmp.from, tmp.to, tmp.type, tmp.clock, tmp.prob);
		if (tmp.level == 380)
			break;
	}
	fclose(fp);
	fclose(fp2);
	return 0;
#endif 
	string current_block_header = "1fdf22ffc2233ff";
	srand(time(NULL));
	int  seed = rand() << 15 | rand();
	int nonce = 0;
	int level = 103;
	LDPC *ptr = new LDPC;
	printf("level : %d, difficulty : %lf\n", level, ptr->get_ldpc_difficulty(level));
	printf("n : %d\t wc : %d\t wr : %d\n", ldpc_level_table[level].n, ldpc_level_table[level].wc, ldpc_level_table[level].wr);
	ptr->set_difficulty(level);
	ptr->initialization();
	ptr->generate_seed(seed);
	ptr->generate_H();
	ptr->generate_Q();
	string current_block_header_with_nonce;

	while (1)
	{
		current_block_header_with_nonce.assign(current_block_header);
		current_block_header_with_nonce += to_string(nonce);
		ptr->generate_hv((unsigned char*)current_block_header_with_nonce.c_str());;
		ptr->decoding();
		if (ptr->decision())
			break;
		nonce++;

	}
	printf("%d\n", nonce);
	ptr->print_word(NULL, 1);
	ptr->print_word("output.txt", 2);
	ptr->print_H("H.txt");

	delete ptr;





}

