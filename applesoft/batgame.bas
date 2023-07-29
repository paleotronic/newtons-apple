  5 GR 
 10 PHYSICS=8192 : GOSUB 5000 : GOSUB 1100
 20 OB=0 : GOSUB 1500 : REM create obj 0
 30 OB=0 : CO=4 : GOSUB 1200 : REM set color to 4
 40 OB=0 : TY=0 : GOSUB 1300 : REM set type to elastic
 50 OB=0 : VX=0 : VY=0 : GOSUB 1400 : REM set velocity (px/sec)
 55 OB=0 : MA=15 : GOSUB 1000 : REM set mass to 15 KG
 60 OB=0 : PX=1 : PY=20 : GOSUB 1600 : REM set object position
 70 OB=0 : WW=1 : HH=5 : GOSUB 2300 : REM set object rect
 80 OB=0 : TY=1 : GOSUB 1300 : REM mechanical

 100 REM ball
 120 OB=1 : GOSUB 1500 : REM create obj 1
 130 OB=1 : CO=7 : GOSUB 1200 : REM set color to 4
 140 OB=1 : TY=0 : GOSUB 1300 : REM set type to elastic
 150 OB=1 : VX=2+INT(RND(1)*3) : VY=2+INT(RND(1)*3) : GOSUB 1400 : REM set velocity (px/sec)
 160 OB=1 : MA=5 : GOSUB 1000 : REM set mass to 15 KG
 170 OB=1 : PX=6 : PY=20 : GOSUB 1600 : REM set object position
 
 200 OB=2:PX=0:PY=0:WW=40:HH=1:CO=3:GOSUB 2400: REM top wall
 210 OB=3:PX=0:PY=39:WW=40:HH=1:CO=3:GOSUB 2400: REM bot wall
 220 OB=4:PX=39:PY=0:WW=1:HH=40:CO=3:GOSUB 2400: REM right wall
 230 GOSUB 500 : REM target

 280 GOSUB 1700 : REM start physics
 281 PP=20 : SC = 0 : LIVES=3
 285 GOSUB 300 : HOME : ? "SCORE: ";SC,"LIVES: ";LIVES
 286 GOSUB 1900 : OB=1 : GOSUB 2200 : IF OO = 1 THEN GOTO 289
 287 FOR Y=1 TO 25 : NEXT Y : GOSUB 400
 288 GOTO 285
 289 LIVES = LIVES - 1 : IF LIVES > 0 THEN GOSUB 600 : GOTO 285
 290 GOSUB 1800 : REM stop physics
 299 END
 
 300 REM move bat
 310 K=PEEK(49152):IF K < 128 THEN RETURN
 320 POKE 49168,0
 330 IF K=234 AND PP > 4 THEN PP = PP - 2 : OB = 0 : PX = 1 : PY = PP : GOSUB 1600
 340 IF K=235 AND PP < 36 THEN PP = PP + 2 : OB = 0 : PX = 1 : PY = PP : GOSUB 1600
 350 RETURN
 
 400 REM check collisions
 410 OB=1 : GOSUB 2500
 420 IF CO=0 THEN RETURN : REM no collision
 430 IF OB=0 THEN SC = SC + 20 : REM bat
 440 IF OB=2 OR OB=3 OR OB=4 THEN SC = SC + 1 : REM wall
 450 IF OB=5 THEN SC = SC + 100 : GOSUB 500 : REM target hit
 460 RETURN
 
 500 REM position target
 520 OB=5 : GOSUB 1500 : REM create obj 5
 530 OB=5 : CO=13 : GOSUB 1200 : REM set color to 13
 540 OB=5 : TY=1 : GOSUB 1300 : REM set type to mechanical
 550 OB=5 : VX=0 : VY=0 : GOSUB 1400 : REM set velocity (px/sec)
 560 OB=5 : PX=20+INT(RND(1)*10)+1 : PY=2+INT(RND(1)*36)+1 : GOSUB 1600 : REM set object position
 570 OB=5 : WW=2 : HH=2 : GOSUB 2300 : REM set object rect
 580 RETURN

 600 HOME : ? "OH NO - TRY AGAIN!" : FOR Y=1 TO 1000: NEXT Y : REM reset play
 605 OB=0 : PP=20 : PX = 1 : PY = PP : GOSUB 1600 : REM reset bat
 610 GOSUB 500
 620 GOSUB 700
 630 RETURN

 700 REM ball
 710 OB=1 : PX=INT(RND(1)*10)+5 : PY=INT(RND(1)*35)+2 : GOSUB 1600
 720 OB=1 : VX=2+INT(RND(1)*3) : VY=2+INT(RND(1)*3) : GOSUB 1400 : REM set velocity (px/sec)
 730 RETURN
 
 1000 REM set-object-mass OB=objectid, MA = mass
 1010 POKE 768,4: POKE 769,OB: POKE 770,MA: CALL PHYSICS
 1020 RETURN
 
 1100 REM init-physics
 1110 POKE 768,0: CALL PHYSICS
 1120 RETURN
 
 1200 REM set-object-color OB=objectid, CO = color
 1210 POKE 768,3: POKE 769,OB: POKE 770,CO: CALL PHYSICS
 1220 RETURN

 1300 REM set-object-type OB=objectid, TY = type (0=elastic, 1=mechanical)
 1310 POKE 768,5: POKE 769,OB: POKE 770,TY: CALL PHYSICS
 1320 RETURN
 
 1400 REM set-object-velocity OB=objectid, VX = xvel, VY = yvel (-127 to +127)
 1405 IF VX < 0 THEN VX = 256 + VX
 1406 IF VY < 0 THEN VY = 256 + VY
 1410 POKE 768,2: POKE 769,OB: POKE 770,VX: POKE 771,VY: CALL PHYSICS
 1420 RETURN
 
 1500 REM define-object OB=object-id
 1510 POKE 768,1: POKE 769,OB: CALL PHYSICS
 1520 RETURN
 
 1600 REM set-object-position OB=objectid, PX = x, PY = y
 1610 POKE 768,6: POKE 769,OB: POKE 770,PX: POKE 771,PY: CALL PHYSICS
 1620 RETURN
 
 1700 REM start-physics
 1710 POKE 768,32: CALL PHYSICS
 1720 RETURN
 
 1800 REM stop-physics
 1810 POKE 768,33: CALL PHYSICS
 1820 RETURN
 
 1900 REM request-update-video
 1910 POKE 768,34: CALL PHYSICS
 1920 RETURN
 
 2000 REM get-object-pos OB=object-id
 2010 POKE 768,7: POKE 769,OB: CALL PHYSICS
 2020 PX = PEEK(768) : PY = PEEK(769)
 2030 RETURN
 
 2100 REM get-object-color OB=object-id
 2110 POKE 768,8: POKE 769,OB: CALL PHYSICS
 2120 CO = PEEK(768) 
 2130 RETURN
 
 2200 REM get-object-oob-state OB=object-id OO=1 (out), OO=0 (in)
 2210 POKE 768,9: POKE 769,OB: CALL PHYSICS
 2220 OO = PEEK(768) 
 2230 RETURN
 
 2300 REM set-object-rect OB=object-id WW=width, HH=HEIGHT
 2310 POKE 768,10: POKE 769,OB: POKE 770,WW: POKE 771,HH: CALL PHYSICS
 2320 RETURN
 
 2400 REM set-block-rect OB=object-id WW=width, HH=HEIGHT, PX=x, PY=y, CO=color
 2410 POKE 768,11: POKE 769,OB: POKE 770,PX: POKE 771,PY: POKE 772,WW: POKE 773, HH: POKE 774, CO: CALL PHYSICS
 2420 RETURN
 
 2500 REM get-object-collisions OB=object-id -> CO=1 (yes) 0 (no) OB=objid
 2510 POKE 768,12: POKE 769,OB: CALL PHYSICS
 2520 CO = PEEK(768) : OB = PEEK(769)
 2530 RETURN
 
 5000 REM load physics driver
 5010 D$=CHR$(4)
 5020 PRINT D$;"BLOAD SERIA"
 5030 RETURN
 
