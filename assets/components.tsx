<?xml version="1.0" encoding="UTF-8"?>
<tileset version="1.2" tiledversion="1.2.1" name="components" tilewidth="16" tileheight="16" tilecount="15" columns="0">
 <grid orientation="orthogonal" width="1" height="1"/>
 <tile id="0">
  <properties>
   <property name="con" value="tb,bt"/>
   <property name="lat" value="250ms"/>
  </properties>
  <image width="16" height="16" source="wire_v.png"/>
 </tile>
 <tile id="1">
  <properties>
   <property name="con" value="tr,rl,lt,rt,lr,tl"/>
   <property name="lat" value="250ms"/>
  </properties>
  <image width="16" height="16" source="wire_t.png"/>
 </tile>
 <tile id="2">
  <properties>
   <property name="con" value="tr,rb,tb,rt,br,bt"/>
   <property name="lat" value="250ms"/>
  </properties>
  <image width="16" height="16" source="wire_r.png"/>
 </tile>
 <tile id="3">
  <properties>
   <property name="con" value="tb,lb,lt,bt,bl,tl"/>
   <property name="lat" value="250ms"/>
  </properties>
  <image width="16" height="16" source="wire_l.png"/>
 </tile>
 <tile id="4">
  <properties>
   <property name="con" value="lr,rl"/>
   <property name="lat" value="250ms"/>
  </properties>
  <image width="16" height="16" source="wire_h.png"/>
 </tile>
 <tile id="5">
  <properties>
   <property name="con" value="tr,tb,tl,rb,rl,bl,rt,bt,lt,br,lr,lb"/>
   <property name="lat" value="250ms"/>
  </properties>
  <image width="16" height="16" source="wire_cross.png"/>
 </tile>
 <tile id="6">
  <properties>
   <property name="con" value="lr,rl,tb,bt"/>
   <property name="lat" value="250ms"/>
  </properties>
  <image width="16" height="16" source="wire_bridge.png"/>
 </tile>
 <tile id="7">
  <properties>
   <property name="con" value="lr,rb,lb,rl,br,bl"/>
   <property name="lat" value="250ms"/>
  </properties>
  <image width="16" height="16" source="wire_b.png"/>
 </tile>
 <tile id="8">
  <properties>
   <property name="con" value="tx,rx,bx,lx"/>
   <property name="lat" value="250ms"/>
   <property name="proc" value="1000"/>
  </properties>
  <image width="16" height="16" source="cpu.png"/>
 </tile>
 <tile id="9">
  <properties>
   <property name="con" value="tx,rx,bx,lx"/>
   <property name="lat" value="250ms"/>
   <property name="sched" value="infinite"/>
  </properties>
  <image width="16" height="16" source="router.png"/>
 </tile>
 <tile id="10">
  <properties>
   <property name="con" value="lb,bl"/>
   <property name="lat" value="250ms"/>
  </properties>
  <image width="16" height="16" source="wire_bl.png"/>
 </tile>
 <tile id="11">
  <properties>
   <property name="con" value="rb,br"/>
   <property name="lat" value="250ms"/>
  </properties>
  <image width="16" height="16" source="wire_br.png"/>
 </tile>
 <tile id="12">
  <properties>
   <property name="con" value="tl,lt"/>
   <property name="lat" value="250ms"/>
  </properties>
  <image width="16" height="16" source="wire_tl.png"/>
 </tile>
 <tile id="13">
  <properties>
   <property name="con" value="rt,tr"/>
   <property name="lat" value="250ms"/>
  </properties>
  <image width="16" height="16" source="wire_tr.png"/>
 </tile>
 <tile id="14">
  <properties>
   <property name="con" value="tx,rx,bx,lx"/>
   <property name="lat" value="250ms"/>
   <property name="sched" value="multitasking"/>
  </properties>
  <image width="16" height="16" source="host.png"/>
 </tile>
</tileset>
