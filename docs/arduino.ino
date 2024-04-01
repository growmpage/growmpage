#include <SoftwareSerial.h>

SoftwareSerial SIM900(7, 8);

void setup() {
  digitalWrite(9, HIGH);
  delay(3000);
  digitalWrite(9, LOW);
  delay(6000);
  SIM900.begin(19200);
  delay(9000);
  SIM900.println("ATD 0049160883447;");
  delay(200);
  SIM900.println();
  delay(9999);
  SIM900.println("ATH");
}

void loop() {
}
