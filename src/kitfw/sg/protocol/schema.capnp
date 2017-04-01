@0xf31f08e525cb6c59;
using Go = import "go.capnp";
$Go.package("protocol");
$Go.import("testpkg");


struct ConcatReplyCapn { 
   retCode  @0:   Int8; 
   val      @1:   Text; 
} 

struct ConcatRequestCapn { 
   userId  @0:   Int64; 
   str1    @1:   Text; 
   str2    @2:   Text; 
} 

struct SumReplyCapn { 
   retCode  @0:   Int8; 
   val      @1:   Int64; 
} 

struct SumRequestCapn { 
   userId  @0:   Int64; 
   num1    @1:   Int64; 
   num2    @2:   Int64; 
} 

##compile with:

##
##
##   capnp compile -ogo ./schema.capnp

