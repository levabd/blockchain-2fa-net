const {createContext, CryptoFactory} = require('sawtooth-sdk/signing')

const context = createContext('secp256k1')
const privateKey = context.newRandomPrivateKey()
const signer = new CryptoFactory(context).newSigner(privateKey)
const crypto = require('crypto')
const _hash = (x) => crypto.createHash('sha512').update(x).digest('hex').toLowerCase()
const FAMILY_NAME = 'tfa';
// const FAMILY_NAME = 'kaztel';
const FAMILY_NAMESPACE = _hash(FAMILY_NAME).substring(0, 6)
const FAMILY_VERSION = '0.1';
const PORT = '8008';
const {createHash} = require('crypto')
const {protobuf} = require('sawtooth-sdk')
const faker = require('faker')
faker.locale = "ru";
const request = require('request')
const fs = require('fs')
const WebSocket = require('ws')

var protobufLib = require('protocol-buffers')

// pass a proto file as a buffer/string or pass a parsed protobuf-schema object
var messages = protobufLib(fs.readFileSync('go/src/tfa/service_client/service_client.proto'))
// var messages = protobufLib(fs.readFileSync('go/src/tfa/service/service.proto'))

const RECORd_NUMBER = 10
let c = 0
const handle = function (transactions, i) {

    const batchHeaderBytes = protobuf.BatchHeader.encode({
        signerPublicKey: signer.getPublicKey().asHex(),
        transactionIds: transactions.map((txn) => txn.headerSignature),
    }).finish()

    const signature1 = signer.sign(batchHeaderBytes)

    const batch = protobuf.Batch.create({
        header: batchHeaderBytes,
        headerSignature: signature1,
        transactions: transactions
    })

    const batchListBytes = protobuf.BatchList.encode({
        batches: [batch]
    }).finish()

    request.post({
        url: `http://127.0.0.1:${PORT}/batches`,
        body: batchListBytes,
        auth: {
            user: 'sawtooth',
            pass: 'z92aGlTdLVYk6mR',
            sendImmediately: true
        },
        headers: {'Content-Type': 'application/octet-stream'}
    }, (err, response) => {
        console.log('err', err);
        c++
        console.log('response', typeof response.body);
        if (response.body['error'] !== undefined) {
            const code = response['error']['code']
            if (code === 31) {
                console.log('make enoter request', err);
                return;
            }
        }

        console.log(response.body)
    })
}


let start = new Date().getTime();
let recordsAdded = 0

// // WebSocket endpoint for the proxy to connect to
// var endpoint = `ws://sawtooth:z92aGlTdLVYk6mR@127.0.0.1:${PORT}/sawtooth-ws/subscriptions`;
// console.log('attempting to connect to WebSocket %j', endpoint);
//
// // initiate the WebSocket connection
// var socket = new WebSocket(endpoint);
//
// socket.on('open', function () {
//     console.log('"open" event!');
//     socket.send(JSON.stringify({
//         'action': 'subscribe',
//         'address_prefixes': [
//             _hash(FAMILY_NAME).substring(0, 6),
//         ]
//     }));
// });
//
// socket.on('message', function (data, flags) {
//     try {
//         const _data = JSON.parse(data)
//         if (_data.state_changes.length) {
//             var end = new Date().getTime();
//             recordsAdded += _data.state_changes.length
//             console.log('8000 mess length', recordsAdded);
//             console.log("8000 Call to onmessage took " + (end - start) + " milliseconds.")
//         } else {
//             console.log('no changes');
//         }
//
//     } catch (e) {
//         console.log('error');
//     }
// });

function getRandomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

const makeFor = (familyName, payload, pn, cntr) => {
    const phoneNumberPart = _hash(pn.toString()).slice(-64)
    let address = _hash(familyName).substring(0, 6) + phoneNumberPart
    const payloadBytes = messages.SCPayload.encode(payload)
    const transactionHeaderBytes = protobuf.TransactionHeader.encode({
        familyName: familyName,
        familyVersion: FAMILY_VERSION,
        inputs: [address],
        outputs: [address],
        signerPublicKey: signer.getPublicKey().asHex(),
        batcherPublicKey: signer.getPublicKey().asHex(),
        dependencies: [],
        payloadSha512: createHash('sha512').update(payloadBytes).digest('hex')
    }).finish()

    const signature0 = signer.sign(transactionHeaderBytes)
    const transaction = protobuf.Transaction.create({
        header: transactionHeaderBytes,
        headerSignature: signature0,
        payload: payloadBytes
    })
    handle([transaction], cntr)
}

var tlist = []
for (let i = 0; i < RECORd_NUMBER; i++) {
    (function (cntr) {

        var pn = '7705' + getRandomInt(999999, 9999999)
        // var pn =  '380934737177'
        // var pn = '77053364711'
        var uin = getRandomInt(99999999999, 999999999999)
        const payload = {
            Action: 0, // create | update | delete
            PhoneNumber: pn,
            PayloadUser: {
                PhoneNumber: pn,
                Uin: `${uin}`,
                Name: faker.name.findName(),
                IsVerified: false,
                Email: `doamin${faker.internet.email()}`,
                Sex: faker.random.arrayElement(['male', 'female']),
                Birthdate: 12452485,
            }
        }
        makeFor('kaztel', payload, pn, cntr)
        makeFor('tfa', payload, pn, cntr)
    })(i);
}
