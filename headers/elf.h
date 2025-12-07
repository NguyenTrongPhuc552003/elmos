// scripts/includes/elf.h
#pragma once
#include <libelf/gelf.h>
 
#define STT_SPARC_REGISTER 3
#define R_386_32 1
#define R_386_PC32 2
#define R_MIPS_HI16 5
#define R_MIPS_LO16 6
#define R_MIPS_26 4
#define R_MIPS_32 2
#define R_ARM_ABS32 2
#define R_ARM_REL32 3
#define R_ARM_PC24 1
#define R_ARM_CALL 28
#define R_ARM_JUMP24 29
#define R_ARM_THM_JUMP24 30
#define R_ARM_THM_PC22 10
#define R_ARM_MOVW_ABS_NC 43
#define R_ARM_MOVT_ABS 44
#define R_ARM_THM_MOVW_ABS_NC 47
#define R_ARM_THM_MOVT_ABS 48
#define R_ARM_THM_JUMP19 51
#define R_AARCH64_ABS64 257
#define R_AARCH64_PREL64 260
