import click
import arrow
from configparser import ConfigParser
from itertools import chain
import os
import random
import requests
import string
import time



def gen_fcfs():
    fcfs = ""
    for i in range(2):
         fcfs += random.choice(string.ascii_letters[0:26]).upper()
    for i in range(2):
        fcfs += random.choice(string.digits[0:9])
    fcfs += "-"
    for i in range(4):
        fcfs += random.choice(string.digits[0:9])
    return fcfs
    
def gen_pi_serial():
    return ''.join([random.choice(string.digits[0:9]).upper() for i in range(10)])

# I couldn't think of a fruit starting with 'd'. Or 'e'. 
fruits = ["apple", "berry", "canteloupe", "daikon", "egg", "fruit", "gorgonzola", "hat", "ingoberry", "jackfruit", "kiwi"]
def gen_tag():
    return ' '.join([random.choice(fruits) for i in range(2)])

# These decorators are applied to `generate`,
# so beware that adding commands modifies the header 
# of that function.
@click.command()
@click.option('--libraries', default=1, help='Number of libraries.')
@click.option('--sensors', default=10, help="Number of sensors.")
@click.option('--weeks', default=1, help="Number of weeks worth of data.")
@click.option('--token', default="nopenopenope", help="User token.")
def generate(libraries, sensors, weeks, token):
    session_id = 0
    patron_ndx = 0
    # Party like it's Y2K! (Was that a Prince hit?)
    # start_date = arrow.get('2000-01-01 00:00:00')
    start_date = arrow.utcnow()
    mfg_ndx = 0 

    for week in range(weeks):
        for dow in range(7):
            session_id += 1
            date = start_date.shift(days=dow).floor('day')
            print(dow)
            for lib in range(libraries):
                fcfs = gen_fcfs()
                tag = gen_tag()        
                for sense in range(sensors):
                    serial = gen_pi_serial()
                    patron_ndx += 1
                    # We only have 20 manufacturers... this has not been
                    # a factor in analysis *yet*``
                    mfg_ndx = int(random.choice(string.digits[0:20]))
                    # `randint` is inclusive of the end of the range... probably.
                    start_hour = random.randint(0, 23)
                    duration = random.randint(0, 12)
                    clamped = duration
                    if start_hour + duration > 23:
                        clamped = 23 - start_hour
                    # print("start hour", start_hour, "duration", duration, "clamped", clamped)
                    start_date = date.shift(hours=start_hour)
                    end_date = start_date.shift(hours=clamped)
                    # print(clamped, start_date, end_date)
                    start_ts = int(time.mktime(start_date.timetuple()))
                    end_ts = int(time.mktime(end_date.timetuple()))
                    
                    if end_ts < start_ts:
                        print(start_ts, end_ts)
                        print("This should never happen.")
                        os.exit()
                    r = requests.post("http://0.0.0.0:8055/items/durations", 
                        headers = {"Authorization" : "Bearer " + token},
                        json = {
                            "pi_serial" : serial,
                            "fcfs_seq_id" : fcfs,
                            "device_tag" : tag,
                            "session_id" : session_id,
                            "patron_index" : patron_ndx,
                            "manufacturer_index" : mfg_ndx,
                            "start" : start_ts,
                            "end" : end_ts
                        })
                    print(r.text)


if __name__ == '__main__':
    generate()