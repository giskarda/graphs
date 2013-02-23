import random
import time
import telnetlib

now = int(time.time())
hourago = now - 3600

tn = telnetlib.Telnet("localhost", 4242)
for ts in xrange(hourago, now, 5):
    rand = random.randint(1,5)
    val = random.randint(1000, 2000)
    host = random.choice(["morte", "male", "distruzione"])
    # put tsd.hbase.flushes 1361637805 0 host=Riccardo-Settis-MacBook-Pro.local
    theput = "put graphs.test.%d %d %d host=%s\n" % (rand, ts, val, host)
    tn.write(theput)

tn.close()
