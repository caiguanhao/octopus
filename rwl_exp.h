/*********************************************************************
 * Octopus Cards Limited
 *
 * File Name            : rwl_exp.h
 * Ported to Linux by   : Vicky Lau (Ported, Nov 2001)
 *                        Kelvin Li (Refined, Nov 2001 - Aug 2002)
 *							3017, ericchan, 13 Sep 13
 * RWL Version (Linux)  : 3017
 * Last Edit            : 13 Sep 13
 *********************************************************************/

#ifndef RWL_EXP_H
#define RWL_EXP_H

#include <stdio.h>

#ifndef _C_
#define _CPP_
#endif

#define AMF_MESSAGE_LEN	100

typedef struct {
	unsigned int    DevID;		/* Device ID */
	unsigned int    OperID;		/* Operator ID */
	unsigned int    DevTime;	/* Device Time */
	unsigned int    CompID;		/* Company ID */
	unsigned int    KeyVer;		/* Key Version */
	unsigned int    EODVer;		/* EOD Version */
	unsigned int    BLVer;		/* Blacklist Version */
	unsigned int    FIRMVer;	/* Firmware Version */
	unsigned int    CCHSVer;	/* CCHS MSG ID */
	unsigned int    CSSer;		/* CS Serial #, Loc ID */
	unsigned int    IntBLVer;	/* Interim Blacklist Version */
	unsigned int	FuncBLVer;	/* Functional Blacklist Version */
} stDevVer;

typedef struct {
	unsigned int	AlertCode;
	unsigned char	ChineseMsgUTF16[AMF_MESSAGE_LEN];
	unsigned char	EnglishMsg[AMF_MESSAGE_LEN];
	unsigned char 	ChineseMsgBig5[AMF_MESSAGE_LEN];
	unsigned char	ChineseMsgUTF8[AMF_MESSAGE_LEN*2];
} stAlertMsgEntry;

//stDevVer 	theVER; //Mantis #7967: [Linux RWL] Blacklist not loaded 

#ifdef _CPP_					
extern "C" {
#endif

extern int Block(unsigned char, unsigned char *);
extern int WriteID(unsigned int);
extern int Deduct(int, unsigned char *, int);
extern int AddValue(int, unsigned char, unsigned char *, int);
extern int AddFund(int, unsigned char, unsigned char *); // 3083 M_Rebate
extern int DrainAR(unsigned char *);
extern int RetrieveUD(char *);
extern int InitComm(unsigned char, unsigned int, unsigned int); // Project 3101 Android Library, Mantis 10510 - item 7
extern int SendCCHS(char *);
extern int Poll(unsigned char, unsigned char, char *);
extern int PollEx(unsigned char, unsigned char, char *); //Mantis #7772
extern int PollExCarPark(unsigned char, unsigned char, char *, unsigned int, int, unsigned char, char *);
extern int PortClose();
extern int Reset();
extern int SendBlack(char *);
extern int DeferDeduct(unsigned char, int, unsigned char *);
extern int GetExtraInfo(unsigned int, unsigned int, unsigned char *);
extern int PIN(int,int);
extern int SendEOD(char *);
extern int FirmUpg(char *);
extern int TxnAmt(int, int, unsigned char, unsigned char);
extern int XFile(char *);
extern int TimeVer(unsigned char *); 
extern int HouseKeeping();
extern int ReadLoyal(unsigned char *);
extern int WriteLoyal(unsigned char*,int);
extern int ReadSPB(unsigned char*, unsigned char*);
extern int WriteSPB(unsigned char*, unsigned char, unsigned char);
extern int InitSPB(unsigned char*);
extern int Format();
extern int InitActionList();
extern int GetSysInfo(unsigned char *bInfo); // End 362 ITRM Blacklist Enhancement
extern int ActionCheck(unsigned int uiCom, unsigned int uiParam, unsigned int *uiResult);
extern int ActionUpdate(unsigned char *ucCom, unsigned char *ucRes);
extern int GetActionResult(unsigned char *ucResult);
extern int InitSubsidyList(void);
extern int CollectSubsidy(unsigned int uiRefNum, unsigned char *ucResult);

#ifdef _CPP_
}
#endif

#endif // RWL_EXP_H


