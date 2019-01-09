inData =
"
#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,dateTime:RFC3339,string,double,double
#group,false,false,true,true,false,true,false,false
#default,,,,,,,,
,result,table,_start,_stop,_time,_measurement,v1,v2
,,0,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:00Z,_m0,0,1
,,0,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:10Z,_m0,1,1
,,0,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:20Z,_m0,2,1
,,0,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:30Z,_m0,3,1
,,0,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:40Z,_m0,4,1
,,0,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:50Z,_m0,5,1
,,1,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:00Z,_m1,0,1
,,1,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:10Z,_m1,2,1
,,1,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:20Z,_m1,4,1
,,1,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:30Z,_m1,6,1
,,1,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:40Z,_m1,8,1
,,1,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:50Z,_m1,6,1
,,1,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:54:00Z,_m1,4,1
,,1,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:54:10Z,_m1,2,1
,,1,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:54:20Z,_m1,0,1
,,2,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:00Z,_m2,0,1
,,2,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:10Z,_m2,8,1
,,2,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:20Z,_m2,2,1
,,2,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:30Z,_m2,6,1
,,3,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:40Z,_m3,1,1
,,3,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:53:50Z,_m3,1,1
,,3,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,2018-05-22T19:54:00Z,_m3,1,1
"
outData =
"#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,string,double,double
#group,false,false,true,true,true,false,false
#default,_result,,,,,,
,result,table,_start,_stop,_measurement,v1,v2
,,0,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,_m0,12.5,5
,,1,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,_m1,32,8
,,2,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,_m2,13,3
,,3,2018-05-22T19:53:00Z,2018-05-22T19:54:50Z,_m3,2,2
"

t_integral_columns = (table=<-) =>
  table |> integral(columns: ["v1", "v2"], unit: 10s)

testingTest(
    name: "integral_columns",
     input: testLoadStorage(csv: inData),
     want: testLoadMem(csv: outData),
     test:  t_integral_columns,
)