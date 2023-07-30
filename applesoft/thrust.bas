 5 GR : BC = 0 : FC = 0 : HOME : ? "J=LEFT, K=RIGHT, I=THRUST, M=BRAKE"
 10 BW=4 : BH=4 : REM block sizes
 20 PHYSICS=8192 : GOSUB 5000 : GOSUB 1100
 30 GOSUB 700 : GOSUB 800 : REM Wall and ship
 35 GOSUB 500 : REM rocks
 40 GOSUB 1700 : REM start physics
 80 FOR Y=1 TO 5 : NEXT Y
 90 GOSUB 1900 : REM update video
 95 GOSUB 300
 110 GOTO 80

 300 REM move ship
 310 K=PEEK(49152):IF K < 128 THEN RETURN
 320 POKE 49168,0
 325 IF K=233 THEN GOSUB 380
 326 IF K=237 THEN GOSUB 390
 330 IF K=234 THEN GOSUB 360
 340 IF K=235 THEN GOSUB 370
 350 RETURN
 360 OB=0 : GOSUB 3200 : HE = HE + 30 : IF HE > 360 THEN HE = HE - 360
 365 OB=0 : GOSUB 3100 : RETURN
 370 OB=0 : GOSUB 3200 : HE = HE - 30 : IF HE < 0 THEN HE = HE + 360
 375 OB=0 : GOSUB 3100 : RETURN
 380 OB=0 : GOSUB 3200 : VE=1 : GOSUB 3000 : REM add velocity
 385 RETURN
 390 OB=0 : GOSUB 3200 : VE=1 : HE=HE-180 : IF HE<0 THEN HE=HE+360:  GOSUB 3000 : REM add velocity
 395 RETURN

 500 REM create ship
 510 OB=0 : GOSUB 1500 : REM create obj 0
 520 OB=0 : EL=100 : GOSUB 2800: REM set elasticity to 20 percent
 530 OB=0 : CO=15 : GOSUB 1200 : REM set color to random
 550 OB=0 : VX=0 : VY=0 : GOSUB 1400 : REM set velocity (px/sec)
 560 OB=0 : PX=20 : PY=20 : GOSUB 1600 : REM set object position
 570 OB=0 : WW=1 : HH=3 : GOSUB 2300 : REM set object rect
 590 RETURN

 700 OB=1:PX=0:PY=0:WW=40:HH=1:CO=0:GOSUB 2400: REM top wall
 710 OB=2:PX=0:PY=39:WW=40:HH=1:CO=0:GOSUB 2400: REM bot wall
 720 OB=3:PX=39:PY=0:WW=1:HH=40:CO=0:GOSUB 2400: REM right wall
 730 OB=4:PX=0:PY=0:WW=1:HH=40:CO=0:GOSUB 2400: REM left wall
 740 RETURN

 800 FOR J=5 TO 15
 810 OB=J: GOSUB 1500 : REM create obj 
 820 OB=J: EL=50 : GOSUB 2800: REM set elasticity to 20 percent
 830 OB=J: CO=INT(RND(1)*10)+5 : GOSUB 1200 : REM set color to random
 850 OB=J: VX=0 : VY=0 : GOSUB 1400 : REM set velocity (px/sec)
 860 OB=J: PX=INT(RND(1)*35)+1 : PY=INT(RND(1)*35)+1 : GOSUB 1600 : REM set object position
 870 OB=J: WW=1+INT(RND(1)*2) : HH=1+INT(RND(1)*3) : GOSUB 2300 : REM set object rect
 875 OB=J: RT=1 : GOSUB 3300 : REM allow rotation
 876 OB=J: HE=INT(RND(1))*360 : GOSUB 3100
 899 NEXT J : RETURN

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

 2600 REM set-velocity-heading OB=object-id VE=vel HE=heading (0-359)
 2605 HH=INT(HE/256) : LL=HE-HH*256
 2610 POKE 768,13: POKE 769,OB: POKE770,VE: POKE771,LL: POKE772,HH: CALL PHYSICS
 2620 RETURN

 2700 REM set-force FO=force HE=heading (0-359)
 2705 HH=INT(HE/256) : LL=HE-HH*256
 2710 POKE 768,14: POKE 769,FO: POKE770,LL: POKE771,HH: CALL PHYSICS
 2720 RETURN

 2800 REM set-object-elasticity OB=object-id, EL=elasticity
 2810 POKE 768,15: POKE 769,OB: POKE 770, EL: CALL PHYSICS
 2820 CO = PEEK(768) 
 2830 RETURN

 2900 REM get-any-oob -> OO=1 (yes) 0 (no) OB=objid
 2910 POKE 768,16: CALL PHYSICS
 2920 OO = PEEK(768) : OB = PEEK(769)
 2930 RETURN

 3000 REM add-velocity-heading OB=object-id VE=vel HE=heading (0-359)
 3005 HH=INT(HE/256) : LL=HE-HH*256
 3010 POKE 768,17: POKE 769,OB: POKE770,VE: POKE771,LL: POKE772,HH: CALL PHYSICS
 3020 RETURN

 3100 REM set-object-heading OB=objectid, HE = heading
 3105 HH=INT(HE/256) : LL=HE-HH*256
 3110 POKE 768,18: POKE 769,OB: POKE 770,LL: POKE771,HH:  CALL PHYSICS
 3120 RETURN

 3200 REM get-object-heading OB=object -> HE = heading
 3210 POKE 768,19: CALL PHYSICS
 3220 HE = PEEK(768) + 256 * PEEK(769)
 3230 RETURN
 
 3300 REM set-object-rotation OB=object -> RT (1) can rotate, (0) cannot rotate
 3310 POKE 768,20: POKE 769,OB: POKE 770,RT: CALL PHYSICS
 3320 RETURN

 5000 REM load physics driver
 5010 D$=CHR$(4)
 5020 PRINT D$;"BLOAD SERIA"
 5030 RETURN
 