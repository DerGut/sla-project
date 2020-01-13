print("Locking db writes");
db.fsyncLock();

print("Waiting for 2min");
sleep(2 * 60 * 1000);

print("Unlocking db writes");
db.fsyncUnlock();
