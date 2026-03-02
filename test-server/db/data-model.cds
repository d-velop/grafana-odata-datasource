namespace test;

entity Temperatures {
    key id          : UUID;
    time            : DateTime;
    sampledAt       : Timestamp;
    measurementDate : Date;
    epoch           : Int64;
    sensorId        : Integer;
    qualityCode     : Int16;
    value1          : Double;
    value2          : Double;
    value3          : Double;
    pressure        : Decimal(7, 2);
    isOutdoor       : Boolean;
    unit            : String;
}

entity Rooms {
    key id : UUID;
    name: String
}
entity TestPrimitives {
    key guid          : UUID;
    dateTimeOffset    : DateTime;
    date              : Date;
    int64             : Int64;
    int32             : Integer;
    int16             : Int16;
    decimal           : Decimal(10, 3);
    double            : Double;
    boolean           : Boolean;
    string            : String;
}
