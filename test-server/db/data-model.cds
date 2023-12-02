namespace test;

entity Temperatures {
    key id : UUID;
    time   : DateTime;
    epoch  : Integer;
    value1 : Double;
    value2 : Double;
    value3 : Double;
}

entity Rooms {
    key id : UUID;
    name: String
}
