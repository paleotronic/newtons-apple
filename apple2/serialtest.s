
          ORG $2000      ; HIRES area unused and safe
COUNTER    = $08         ; Rotating counter for RXbyte
BUFFER     = $09         ; Where IRQ will store RXbyte
SLOT       = $20         ; This is slot 2 where SSC is
DATAREG     = $C088+SLOT
STATUSREG   = $C089+SLOT
COMMANDREG  = $C08A+SLOT
CONTROLREG  = $C08B+SLOT
HOME       = $FC58       ; CLS & Cursor to top POS
COUT       = $FDED       ; Like C64 $FFD2
KYBD       = $C000       ; Keyboard map point
STROBE     = $C010       ; Keyboard strobe map point
HEXOUT     = $FDDA       ; ROM call to write out hex
LINEBREAK  = $FD8E

CMD_LEN    = $CD
CMD_LO     = $CE
CMD_HI     = $CF

MLICMD     = $300
MLIARGS    = $301

CMD_INIT   = 0
CMD_CREATE = 1
CMD_SETVEL = 2
CMD_SETCOL = 3
CMD_SETMASS = 4
CMD_SETTYPE = 5
CMD_SETPOS  = 6
CMD_START   = 32
CMD_STOP    = 33
CMD_REQDIFF = 34
CMD_GETPOS  = 7
CMD_GETCOL  = 8
CMD_GETOOB  = 9
CMD_SETRECT = 10
CMD_BLKRECT = 11
CMD_GETCOLL = 12
CMD_SETVELH = 13
CMD_SETFORCE = 14
CMD_SETELAST = 15
CMD_ALLOOB = 16
CMD_ADDVELH = 17
CMD_SETH = 18
CMD_GETH = 19
CMD_SETSPIN = 20

ENTRYPOINT
        ; this is where user CALL()'s come in... 
        ; $FA contains a command
        JMP ACTUALENTRY
CONTROLVAL DB %00011111        
ACTUALENTRY
        LDA MLICMD
        CMP #CMD_INIT
        BEQ JP_INIT
        CMP #CMD_CREATE
        BEQ JP_CREATE
        CMP #CMD_SETVEL
        BEQ JP_SETVEL
        CMP #CMD_SETCOL
        BEQ JP_SETCOL
        CMP #CMD_SETMASS
        BEQ JP_SETMASS
        CMP #CMD_SETTYPE
        BEQ JP_SETTYPE
        CMP #CMD_SETPOS
        BEQ JP_SETPOS
        CMP #CMD_START
        BEQ JP_START
        CMP #CMD_STOP
        BEQ JP_STOP
        CMP #CMD_REQDIFF
        BEQ JP_REQDIFF
        CMP #CMD_GETPOS
        BEQ JP_GETPOS
        CMP #CMD_GETCOL
        BEQ JP_GETCOL
        CMP #CMD_GETOOB
        BEQ JP_GETOOB
        CMP #CMD_SETRECT
        BEQ JP_SETRECT
        CMP #CMD_BLKRECT
        BEQ JP_BLKRECT
        CMP #CMD_GETCOLL
        BEQ JP_GETCOLL
        CMP #CMD_SETVELH
        BEQ JP_SETVELH
        CMP #CMD_SETFORCE
        BEQ JP_SETFORCE
        CMP #CMD_SETELAST
        BEQ JP_SETELAST
        CMP #CMD_ALLOOB
        BEQ JP_ALLOOB
        CMP #CMD_ADDVELH
        BEQ JP_ADDVELH
        CMP #CMD_SETH
        BEQ JP_SETH
        CMP #CMD_GETH
        BEQ JP_GETH
        CMP #CMD_SETSPIN
        BEQ JP_SETSPIN
        RTS

JP_INIT 
        JMP P_INIT
JP_CREATE
        JMP P_CREATE
JP_SETVEL
        JMP P_SETVEL
JP_SETCOL
        JMP P_SETCOL
JP_SETMASS
        JMP P_SETMASS
JP_SETTYPE
        JMP P_SETTYPE
JP_SETPOS
        JMP P_SETPOS
JP_START
        JMP P_START
JP_STOP 
        JMP P_STOP
JP_REQDIFF
        JMP P_REQDIFF
JP_GETPOS
        JMP P_GETPOS
JP_GETCOL
        JMP P_GETCOL
JP_GETOOB
        JMP P_GETOOB
JP_SETRECT
        JMP P_SETRECT
JP_BLKRECT
        JMP P_BLKRECT
JP_GETCOLL
        JMP P_GETCOLL
JP_SETVELH
        JMP P_SETVELH
JP_SETFORCE
        JMP P_SETFORCE
JP_SETELAST
        JMP P_SETELAST
JP_ALLOOB
        JMP P_ALLOOB
JP_ADDVELH
        JMP P_ADDVELH
JP_SETH
        JMP P_SETH
JP_GETH
        JMP P_GETH
JP_SETSPIN
        JMP P_SETSPIN

P_INIT
        JSR INIT
        LDX #4
        LDA #<INITPHYSICS
        LDY #>INITPHYSICS
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        RTS

P_START
        LDX #START_L
        LDA #<START
        LDY #>START 
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS 

P_STOP
        LDX #STOP_L
        LDA #<STOP
        LDY #>STOP 
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS 

P_ALLOOB
        LDX #ALLOOB_L
        LDA #<ALLOOB
        LDY #>ALLOOB 
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        LDA COMMANDBUFFER+4
        STA MLICMD+1
        RTS 

P_REQDIFF
        LDX #REQDIFF_L
        LDA #<REQDIFF
        LDY #>REQDIFF
        JSR SENDCOMMAND
        JSR RECVCOMMAND
REQDIFFCHECK
        LDA COMMANDBUFFER
        CMP #$85
        BEQ MEMUPDATEMO
        CMP #$81
        BEQ MEMUPDATE
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS 

MEMUPDATEMO
        JSR MEMWRITE  ; write the data
        LDX #REQUESTMORE_L
        LDA #<REQUESTMORE
        LDY #>REQUESTMORE
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        JMP REQDIFFCHECK
MEMUPDATE
        JMP MEMWRITE

P_GETCOLL
        LDA MLIARGS
        STA GETCOLL0 ; object number
        LDX #GETCOLL_L
        LDA #<GETCOLL
        LDY #>GETCOLL
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        LDA COMMANDBUFFER+4
        STA MLICMD+1
        RTS 

P_GETPOS
        LDA MLIARGS
        STA GETPOS0 ; object number
        LDX #GETPOS_L
        LDA #<GETPOS
        LDY #>GETPOS
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        LDA COMMANDBUFFER+4
        STA MLICMD+1
        RTS 

P_GETCOL
        LDA MLIARGS
        STA GETCOL0 ; object number
        LDX #GETCOL_L
        LDA #<GETCOL
        LDY #>GETCOL
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS 

P_GETH
        LDA MLIARGS
        STA GETH0 ; object number
        LDX #GETH_L
        LDA #<GETH
        LDY #>GETH
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        LDA COMMANDBUFFER+4
        STA MLICMD+1
        RTS 

P_GETOOB
        LDA MLIARGS
        STA GETOOB0 ; object number
        LDX #GETOOB_L
        LDA #<GETOOB
        LDY #>GETOOB
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS 

P_CREATE
        LDA MLIARGS
        STA CREATEOBJ0 ; object number
        LDX #CREATEOBJ_L
        LDA #<CREATEOBJ
        LDY #>CREATEOBJ
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS     

P_SETVEL
        LDA MLIARGS
        STA VELOCITY0 ; object number
        LDA MLIARGS+1
        STA VELOCITY1 ; vel x 
        LDA MLIARGS+2
        STA VELOCITY2
        LDX #VELOCITY_L
        LDA #<VELOCITY
        LDY #>VELOCITY
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS     

P_SETVELH
        LDA MLIARGS
        STA VELOCITYH0 ; object number
        LDA MLIARGS+1
        STA VELOCITYH1 ; vel  
        LDA MLIARGS+2
        STA VELOCITYH2 ; heading lo
        LDA MLIARGS+3
        STA VELOCITYH3 ; heading hi
        LDX #VELOCITYH_L
        LDA #<VELOCITYH
        LDY #>VELOCITYH
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS 

P_ADDVELH
        LDA MLIARGS
        STA ADDVELH0 ; object number
        LDA MLIARGS+1
        STA ADDVELH1 ; vel  
        LDA MLIARGS+2
        STA ADDVELH2 ; heading lo
        LDA MLIARGS+3
        STA ADDVELH3 ; heading hi
        LDX #ADDVELH_L
        LDA #<ADDVELH
        LDY #>ADDVELH
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS 

P_SETFORCE
        LDA MLIARGS
        STA FORCE0 ; strength
        LDA MLIARGS+1
        STA FORCE1 ; heading lo
        LDA MLIARGS+2
        STA FORCE2 ; heading hi
        LDX #FORCE_L
        LDA #<FORCE
        LDY #>FORCE
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS 

P_SETRECT
        LDA MLIARGS
        STA SETRECT0 ; object number
        LDA MLIARGS+1
        STA SETRECTW ; x1
        LDA MLIARGS+2
        STA SETRECTH ; y1
        LDA MLIARGS+3
        LDX #SETRECT_L
        LDA #<SETRECT
        LDY #>SETRECT
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS  

P_BLKRECT
        LDA MLIARGS
        STA BLKRECT0 ; object number
        LDA MLIARGS+1
        STA BLKRECTX ; x
        LDA MLIARGS+2
        STA BLKRECTY ; y
        LDA MLIARGS+3
        STA BLKRECTW ; w
        LDA MLIARGS+4
        STA BLKRECTH ; h
        LDA MLIARGS+5
        STA BLKRECTC ; h
        LDX #BLKRECT_L
        LDA #<BLKRECT
        LDY #>BLKRECT
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS  

P_SETCOL
        LDA MLIARGS
        STA COLOR0 ; object number
        LDA MLIARGS+1
        AND #$0f
        STA COLOR1 ; color
        LDX #COLOR_L
        LDA #<COLOR
        LDY #>COLOR
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS   

P_SETH
        LDA MLIARGS
        STA SETH0 ; object number
        LDA MLIARGS+1
        STA SETH1 ; angle lo
        LDA MLIARGS+2
        STA SETH2 ; angle lo
        LDX #SETH_L
        LDA #<SETH
        LDY #>SETH
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS   

P_SETELAST
        LDA MLIARGS
        STA ELASTIC0 ; object number
        LDA MLIARGS+1
        STA ELASTIC1 ; color
        LDX #ELASTIC_L
        LDA #<ELASTIC
        LDY #>ELASTIC
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS   

P_SETMASS
        LDA MLIARGS
        STA MASS0 ; object number
        LDA MLIARGS+1
        STA MASS1 ; mass
        LDX #MASS_L
        LDA #<MASS
        LDY #>MASS
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        * LDA COMMANDBUFFER+3
        * STA MLICMD
        RTS   

P_SETTYPE
        LDA MLIARGS
        STA TYPE0 ; object number
        LDA MLIARGS+1
        STA TYPE1 ; type
        LDX #TYPE_L
        LDA #<TYPE
        LDY #>TYPE
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS 

P_SETSPIN
        LDA MLIARGS
        STA SPIN0 ; object number
        LDA MLIARGS+1
        STA SPIN1 ; true/false
        LDX #SPIN_L
        LDA #<SPIN
        LDY #>SPIN
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS 

P_SETPOS
        LDA MLIARGS
        STA POS0 ; object number
        LDA MLIARGS+1
        STA POS1 ; x
        LDA MLIARGS+2
        STA POS2 ; y
        LDX #POS_L
        LDA #<POS
        LDY #>POS
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        LDA COMMANDBUFFER+3
        STA MLICMD
        RTS 

************** INIT Sets up the SSC and IRQ stuff
INIT
          LDA #$00
          STA STATUSREG
          STA COUNTER    ; set counter to 0
          STA BUFFSIZE

********* LOAD VALUES INTO CONTROLREG
* BIT [7]    0 - ONE STOP BIT / 1 TWO STOP BITS
* BIT [6-5]  LENGTH - 8 (00) / 7 (01) / 6 (10) / 5 (11)
* BIT [4]    CLK - EXTERNAL (0) / INTERNAL (1)
* BIT [3-0]  BAUD - 9600 (1110) / 19.2K (1111)
*                   4800 (1100) / 2400  (1010)
******* USING N81 INTERNAL AT 9600 BAUD
          LDA CONTROLVAL
          STA CONTROLREG

********* LOAD VALUES INTO COMMANDREG
* BIT [7-5] PARITY- (000)DIS  (001)O  (011)E  (101)M (111)S
* BIT [4]   ECHO -  (0) NORM / (1) ON
* BIT [3-2] TXINT- (00)DIS:RTSH (01)EN:RTSL
*                  (10)DIS:RTSL (11)DIS:SBRK
* BIT [1] RXINT - 0 ENABLED / 1 DISBLED
* BIT [0] DTR - 0 DISABLE RXTX,DTR HIGH / 1 ENABLE RXTX DTR LOW
********* USING RX INTERRUPTS ENABLED DTR LOW ENABLE RXTX
          LDA #%00001001
          STA COMMANDREG
          LDA DATAREG    ; Pull from DATAREG on init just because
          ;
        LDX #HELLO_L
        LDA #<HELLO
        LDY #>HELLO
        JSR SENDCOMMAND
        JSR RECVCOMMAND
        RTS
    
MODEMINIT  ASC 'ATZ'
           DB 13
           ASC 'ATDT407'
           DB 13

HELLO_L = 5
HELLO
           DB $FF,$2,$0,$AD,$DE

INITPHYSICS
           DB $01,$1,$0,$0

REQUESTMORE_L = 3
REQUESTMORE
           DB $40,$0,$0

CREATEOBJ_L = 4
CREATEOBJ
           DB $02 ; command byte
           DB $01,$00 ; size
CREATEOBJ0 DB $00 ; object num

FORCE_L = 6
FORCE
           DB $09 ; command byte
           DB $03,$00 ; size
FORCE0  DB $00 ; force
FORCE1  DB $00 ; heading lo
FORCE2  DB $00 ; heading hi

VELOCITY_L = 6
VELOCITY
           DB $05 ; command byte
           DB $03,$00 ; size
VELOCITY0  DB $00 ; object num
VELOCITY1  DB $00 ; vel x  int8
VELOCITY2  DB $00 ; vel y  int8

VELOCITYH_L = 7
VELOCITYH
           DB $18 ; command byte
           DB $04,$00 ; size
VELOCITYH0  DB $00 ; object num
VELOCITYH1  DB $00 ; vel (0-255)
VELOCITYH2  DB $00 ; heading lo
VELOCITYH3  DB $00 ; heading hi

ADDVELH_L = 7
ADDVELH
           DB $20 ; command byte
           DB $04,$00 ; size
ADDVELH0  DB $00 ; object num
ADDVELH1  DB $00 ; vel (0-255)
ADDVELH2  DB $00 ; heading lo
ADDVELH3  DB $00 ; heading hi

COLOR_L = 5
COLOR
           DB $0f ; command byte
           DB $02,$00 ; size
COLOR0  DB $00 ; object num
COLOR1  DB $00 ; col (0-15)

SETH_L = 5
SETH
           DB $21 ; command byte
           DB $02,$00 ; size
SETH0  DB $00 ; object num
SETH1  DB $00 ; col (0-15)
SETH2  DB $00

ELASTIC_L = 5
ELASTIC
           DB $08 ; command byte
           DB $02,$00 ; size
ELASTIC0   DB $00 ; object num
ELASTIC1   DB $00 ; elasticity (0-100)

MASS_L = 5
MASS
           DB $04 ; command byte
           DB $02,$00 ; size
MASS0  DB $00 ; object num
MASS1  DB $00 ; mass 0-255 KG

TYPE_L = 5
TYPE
           DB $11 ; command byte
           DB $02,$00 ; size
TYPE0  DB $00 ; object num
TYPE1  DB $00 ; type: 0 = elastic, 1 = mechanical

SPIN_L = 5
SPIN
           DB $23 ; command byte
           DB $02,$00 ; size
SPIN0      DB $00 ; object num
SPIN1      DB $00 ; 1 = enable, 0 = disable

POS_L = 6
POS
           DB $06 ; command byte
           DB $03,$00 ; size
POS0       DB $00 ; object num
POS1       DB $00 ; x
POS2       DB $00 ; y

START_L = 3
START
           DB $12
           DB $00,$00

STOP_L = 3
STOP
           DB $13
           DB $00,$00

ALLOOB_L = 3
ALLOOB
           DB $19
           DB $00,$00

REQDIFF_L = 3
REQDIFF
           DB $10
           DB $00,$00

GETPOS_L = 4
GETPOS
           DB $14
           DB $01,$00
GETPOS0    DB $00

GETCOLL_L = 4
GETCOLL
           DB $17
           DB $01,$00
GETCOLL0   DB $00

GETCOL_L = 4
GETCOL
           DB $15
           DB $01,$00
GETCOL0    DB $00

GETH_L = 4
GETH
           DB $22
           DB $01,$00
GETH0      DB $00

GETOOB_L = 4
GETOOB
           DB $16
           DB $01,$00
GETOOB0    DB $00

SETRECT_L = 6
SETRECT
           DB $0d
           DB $03,$00
SETRECT0   DB $00
SETRECTW   DB $00
SETRECTH   DB $00

BLKRECT_L = 9
BLKRECT
           DB $0b
           DB $06,$00
BLKRECT0   DB $00
BLKRECTX   DB $00
BLKRECTY   DB $00
BLKRECTW   DB $00
BLKRECTH   DB $00
BLKRECTC   DB $0f

************** SEND Check status and send a byte to ACIA
SEND
          TAX
:CHECK
          LDA STATUSREG  ; Load status register
          AND #$10       ; mask for ready bit
          BEQ :CHECK     ; not ready? keep checking
          STX DATAREG    ; ready?? store byte
          RTS


SENDCOMMAND             ; x contains count of bytes, a = low of command address, y = high of command address
        STX CMD_LEN
        STA CMD_LO
        STY CMD_HI
        LDY #0
:SENDLOOP
        LDA (CMD_LO),Y
        JSR SEND
        INY 
        CPY CMD_LEN
        BNE :SENDLOOP
        RTS   

RECVCOMMAND
        LDA #0
        STA BUFFSIZE
        STA BUFFCOMPLETE
:LOOPRECV
        JSR CHECKBUFFER
        LDA BUFFCOMPLETE
        CLC
        CMP #1
        BCC :LOOPRECV
        RTS

; command handler
CHECKBUFFER
          LDA STATUSREG
          AND #8
          BNE HASDATA
          RTS
HASDATA
          LDX BUFFSIZE
          LDA DATAREG
          CMP #$7F
          BNE CHECKPRE
          LDY #1
          STY PREAMBLE
CHECKPRE
          LDY PREAMBLE
          CPY #1
          BEQ CONTINUE
          RTS
CONTINUE
          STA COMMANDBUFFER,X
          INX
          STX BUFFSIZE
          CPX #3
          BCS CHECKSIZE
          LDA #0
          STA BUFFCOMPLETE
          RTS      
CHECKSIZE
          LDX COMMANDBUFFER+1
          INX
          INX
          INX
          CPX BUFFSIZE
          BEQ GOTBUFFEROK
          LDA #0
          STA BUFFCOMPLETE
          RTS
GOTBUFFEROK
          LDA #1
          STA BUFFCOMPLETE
          RTS

; commands here
MEMCOUNT 
        DB 0 ; acts as an update counter
MEMWRITE
        LDX #0
        LDY #0
        LDA COMMANDBUFFER+3
        STA MEMCOUNT
:MEMWRITELP
        LDA COMMANDBUFFER+4,X
        STA MEMZPADDR
        LDA COMMANDBUFFER+5,X
        STA MEMZPADDR+1    
        LDA COMMANDBUFFER+6,X   
        STA (MEMZPADDR),Y
        INX
        INX
        INX
        DEC MEMCOUNT
        BNE :MEMWRITELP
:MEMWRITEDN
        LDA #1
        STA MLICMD
        LDA #0
        STA BUFFSIZE
        RTS
; memory fill
MEMFILL
        LDA MEMFILLBEG
        STA MEMZPADDR
        LDA MEMFILLBEG+1
        STA MEMZPADDR+1
        LDA MEMFILLVAL          ; value to use  
        LDY MEMFILLCNTLO
        LDX MEMFILLCNTHI 
MEMFILLLOOP
        STA (MEMZPADDR),Y
        DEY
        CPY #$ff
        BNE MEMFILLLOOP
        INC MEMZPADDR+1
        DEX
        CPX #255
        BNE MEMFILLLOOP
        RTS

PREAMBLE DB 0
BUFFSIZE DB 0
BUFFCOMPLETE DB 0
COMMANDBUFFER
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0

SENDCOUNT DB 0
SENDBUFFER 
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
          DB 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0

MEMFILLBEG = COMMANDBUFFER+3
MEMZPADDR = $FA
MEMFILLCNTLO = COMMANDBUFFER+5
MEMFILLCNTHI = COMMANDBUFFER+6
MEMFILLVAL = COMMANDBUFFER+7
