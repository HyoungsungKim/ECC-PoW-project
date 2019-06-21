from threading import Thread, Event

import time
from flask import Flask, jsonify, request
import requests
import hashlib
import json
from urllib.parse import urlparse
from uuid import uuid4

import random

from blockchain_core import blockchain
from blockchain_core import app
from blockchain_core import INTERRUPT_EVENT1
from blockchain_core import INTERRUPT_EVENT2
from blockchain_core import STOP_EVENT

trading_addresses=[
    'Node1',
    'Node2',
    'Node3'
]

class LoopThread(Thread):
    def __init__(self, stop_event, interrupt_event1, interrupt_event2):
        self.stop_event = stop_event
        self.interrupt_event1 = interrupt_event1
        self.interrupt_event2 = interrupt_event2
        Thread.__init__(self)

    def run(self):
        while len(blockchain.chain) < 20 and not self.stop_event.is_set():
            self.loop_process()
            if self.interrupt_event1.is_set():
                self.interrupted_process1()
                self.interrupt_event1.clear()
            if self.interrupt_event2.is_set():
                self.interrupted_process2()
                self.interrupt_event2.clear()

    def loop_process(self):
        TX_period = 20
        senderInd = random.randrange(0,3)
        sender=trading_addresses[senderInd]
        recipientInd=random.randrange(0,3)
        recipient=trading_addresses[recipientInd]
        amount = random.randrange(1, 100)
        blockchain.new_transaction(sender, recipient, amount)
        blockchain.awaiting_transactions.append(blockchain.current_transactions[-1])
        blockchain.announcement()
        print("Transaction generated")
        time.sleep(TX_period)

    def interrupted_process1(self):
        blockchain.resolve_conflicts()
        blockchain.update_awaiting_TX()
        print("Transactions in awaiting are updated")

    def interrupted_process2(self):
        #get_transactions
        print("Transaction transfer")
        #time.sleep(1)

thread = LoopThread(STOP_EVENT, INTERRUPT_EVENT1, INTERRUPT_EVENT2)

@app.route('/transactions/auto', methods=['get'])
def auto_transaction():
    start_time = time.time()
    thread.start()
    end_time = time.time()
    time_consumption = end_time - start_time
    response = {
        'message': 'Auto-TX completed',
        'TX length': len(blockchain.current_transactions),
        'starting time': start_time,
        'ending time': end_time,
        'time consumption': time_consumption,
    }
    return jsonify(response), 200

if __name__ == '__main__':
    from argparse import ArgumentParser

    parser = ArgumentParser()
    parser.add_argument('-p', '--port', default=5000, type=int, help='port to listen on')
    args = parser.parse_args()
    port = args.port

    app.run(host='127.0.0.24', port=2000)