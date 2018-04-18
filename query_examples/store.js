const {createContext, CryptoFactory} = require('sawtooth-sdk/signing')

const context = createContext('secp256k1')
const privateKey = context.newRandomPrivateKey()
const signer = new CryptoFactory(context).newSigner(privateKey)
const crypto = require('crypto')

const _hash = (x) => crypto.createHash('sha512').update(x).digest('hex').toLowerCase()
const cbor = require('cbor')
const FAMILY_NAME = 'tfa';
// const FAMILY_NAME = 'kaztel';
const FAMILY_NAMESPACE = _hash(FAMILY_NAME).substring(0, 6)
const FAMILY_VERSION = '0.1';
const PORT = '8008';
const {createHash} = require('crypto')
const {protobuf } = require( 'sawtooth-sdk')
const faker = require('faker')
faker.locale = "ru";
const request = require('request')
const fs = require('fs')
const WebSocket = require('ws')

var protobufLib = require('protocol-buffers')

// pass a proto file as a buffer/string or pass a parsed protobuf-schema object
// var messages = protobufLib(fs.readFileSync('go/src/tfa/service_client/service_client.proto'))
var messages = protobufLib(fs.readFileSync('go/src/tfa/service/service.proto'))

const RECORd_NUMBER = 1
let c = 0
let e = 0
const makeRequest = (data) => {
    request.post({
        url: `http://127.0.0.1:${PORT}/batches`,
        body: batchListBytes,
        headers: {'Content-Type': 'application/octet-stream'}
    }, (err, response) => {
        c++
        console.log('c', c);
        if (i === RECORd_NUMBER) {
            console.log(response.body)
        }
    }, err => {
        console.log('error', err);

    })
}
// let start = new Date().getTime();

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

    const APIs = [
        '8000', '8001', '8002'
    ];

    const randomPort = APIs[Math.floor(Math.random() * APIs.length)];
    request.post({
        url: `http://127.0.0.1:${PORT}/sawtooth/batches`,
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
                makeRequest(batchListBytes)
                return;
            }
        }

        console.log(response.body)
        // var end = new Date().getTime();
        // console.log("8000 Call to onmessage took " + (end - start) + " milliseconds.")
    })
}

// let ws = new WebSocket(`ws:127.0.0.1:${PORT}/subscriptions`)
// ws.onopen = () => {
//     ws.send(JSON.stringify({
//         'action': 'subscribe',
//         'address_prefixes': [
//             _hash(FAMILY_NAME).substring(0, 6),
//         ]
//     }));
// }
// let start = new Date().getTime();
// let recordsAdded = 0
// ws.onmessage = (mess) => {
//     try {
//         const data = JSON.parse(mess.data)
//         if (data.state_changes.length) {
//             var end = new Date().getTime();
//             recordsAdded+=data.state_changes.length
//             console.log('8000 mess length', recordsAdded);
//             console.log("8000 Call to onmessage took " + (end - start) + " milliseconds.")
//         } else {
//             console.log('no changes');
//         }
//
//     } catch (e) {
//         console.log('error'), e;
//     }
// }
//
// ws.onclose = () => {
//     ws.send(JSON.stringify({'action': 'unsubscribe'}));
// }

function getRandomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

var tlist = []
for (let i = 0; i <= RECORd_NUMBER; i++) {
    (function (cntr) {

         //var pn = '7705' + getRandomInt(999999, 9999999)
        // var pn =  '77053237001'
        var pn =  '77059127941'
        var uin = getRandomInt(99999999999, 999999999999)
        const payload = {
            Action: 0, // create | update | delete
            PhoneNumber: pn,
            PayloadUser: {
                PhoneNumber: pn,
                Uin: uin,
                Name: faker.name.findName(),
                IsVerified: false,
                Email: `doamin${faker.internet.email()}`,
                Sex: faker.random.arrayElement(['male', 'female']),
                Birthdate: 12452485,
            }
        }
        const phoneNumberPart = _hash(pn.toString()).slice(-64)

        // let address = FAMILY_NAMESPACE + _hash(payload.User.Uin +payload.User.PhoneNumber).slice(-64)
        let address = FAMILY_NAMESPACE + phoneNumberPart

        // console.log('address', address);

        const payloadBytes = messages.SCPayload.encode(payload)
        console.log('payloadBytes.length.', payloadBytes.length);

        // const payloadBytes = cbor.encode(payload)

        const transactionHeaderBytes = protobuf.TransactionHeader.encode({
            familyName: FAMILY_NAME,
            familyVersion: FAMILY_VERSION,
            inputs: [address],
            outputs: [address],
            signerPublicKey: signer.getPublicKey().asHex(),
            // In this example, we're signing the batch with the same private key,
            // but the batch can be signed by another party, in which case, the
            // public key will need to be associated with that key.
            batcherPublicKey: signer.getPublicKey().asHex(),
            // In this example, there are no dependencies.  This list should include
            // an previous transaction header signatures that must be applied for
            // this transaction to successfully commit.
            // For example,
            // dependencies: ['540a6803971d1880ec73a96cb97815a95d374cbad5d865925e5aa0432fcf1931539afe10310c122c5eaae15df61236079abbf4f258889359c4d175516934484a'],
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

         // if (tlist.length === 30) {
         //     handle(tlist, cntr)
         //     tlist = []
         // } else {
         //     tlist.push(transaction)
         // }
    })(i);
}
