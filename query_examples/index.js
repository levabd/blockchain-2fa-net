const {createContext, CryptoFactory} = require('sawtooth-sdk/signing')

const context = createContext('secp256k1')
const privateKey = context.newRandomPrivateKey()
const signer = new CryptoFactory(context).newSigner(privateKey)
const crypto = require('crypto')

const _hash = (x) => crypto.createHash('sha512').update(x).digest('hex').toLowerCase()
const cbor = require('cbor')
const FAMILY_NAME = 'tfa';
const FAMILY_NAMESPACE = _hash(FAMILY_NAME).substring(0, 6)
const FAMILY_VERSION = '0.1';
const {createHash} = require('crypto')
const {protobuf} = require('sawtooth-sdk')
const faker = require('faker')
faker.locale = "ru";
const request = require('request')
const WebSocket = require('ws')

const RECORd_NUMBER = 100
let c = 0
let e = 0
const makeRequest = (data) => {
    request.post({
        url: `http://127.0.0.1:8008/batches`,
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
        url: `http://127.0.0.1:8008/batches`,
        body: batchListBytes,
        headers: {'Content-Type': 'application/octet-stream'}
    }, (err, response) => {

        c++
        console.log('response', typeof response.body);
        if (response.body['error']!==undefined) {
            const code = response['error']['code']
            if (code === 31) {
                console.log('make enoter request', err);
                makeRequest(batchListBytes)
                return;
            }
        }

        console.log('c', c);
        if (c === RECORd_NUMBER) {
            console.log(response.body)
        }
    })
}


let ws = new WebSocket(`ws:127.0.0.1:8008/subscriptions`)
ws.onopen = () => {
    ws.send(JSON.stringify({
        'action': 'subscribe',
        'address_prefixes': [
            _hash('kaztel').substring(0, 6),
        ]
    }));
}
var start = new Date().getTime();
ws.onmessage = (mess) => {
    try {
        const data = JSON.parse(mess.data)
        if (data.state_changes) {
            var end = new Date().getTime();
            console.log('8000 mess length', data.state_changes.length);
            console.log("8000 Call to onmessage took " + (end - start) + " milliseconds.")
        } else {
            console.log('no changes');
        }

    } catch (e) {
        console.log('error'), e;
    }
}

ws.onclose = () => {
    ws.send(JSON.stringify({'action': 'unsubscribe'}));
}

// let ws1 = new WebSocket(`ws:127.0.0.1:8001/subscriptions`)
// ws1.onopen = () => {
//     ws1.send(JSON.stringify({
//         'action': 'subscribe',
//         'address_prefixes': [
//             _hash('kaztel').substring(0, 6),
//         ]
//     }));
// }
// var start1 = new Date().getTime();
// ws1.onmessage = (mess) => {
//     try {
//         const data = JSON.parse(mess.data)
//         if (data.state_changes){
//             var end = new Date().getTime();
//             console.log('8001 mess length',data.state_changes.length);
//             console.log("8001 Call to onmessage took " + (end - start1) + " milliseconds.")
//         } else{
//             console.log('no changes');
//         }
//
//     } catch(e){
//         console.log('error'), e;
//     }
// }
//
// ws1.onclose = () => {ws1.send(JSON.stringify({'action': 'unsubscribe'}));}
//
// let ws2 = new WebSocket(`ws:127.0.0.1:8002/subscriptions`)
// ws2.onopen = () => {
//     ws2.send(JSON.stringify({
//         'action': 'subscribe',
//         'address_prefixes': [
//             _hash('kaztel').substring(0, 6),
//         ]
//     }));
// }
// var start2 = new Date().getTime();
// ws2.onmessage = (mess) => {
//     try {
//         const data = JSON.parse(mess.data)
//         if (data.state_changes){
//             var end = new Date().getTime();
//             console.log('8002 mess length',data.state_changes.length);
//             console.log("8002 Call to onmessage took " + (end - start2) + " milliseconds.")
//         } else{
//             console.log('no changes');
//         }
//
//     } catch(e){
//         console.log('error'), e;
//     }
// }
//
// ws2.onclose = () => {
//     ws2.send(JSON.stringify({
//         'action': 'unsubscribe'
//     }));
// }

function getRandomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

var tlist = []
for (let i = 0; i <= RECORd_NUMBER; i++) {
    (function (cntr) {

        var pn = '7705' + getRandomInt(999999, 9999999)
        var uin = getRandomInt(99999999999, 999999999999)

        const payload = {
            Action: 'create', // create | update | delete
            PhoneNumber: pn,
            User: {
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

        console.log('address', address);

        const payloadBytes = cbor.encode(payload)

        const transactionHeaderBytes = protobuf.TransactionHeader.encode({
            familyName: 'tfa',
            familyVersion: '0.1',
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

        // if (tlist.length===30){
        //     handle(tlist, cntr)
        //     tlist=[]
        // } else {
        //     tlist.push(transaction)
        // }

        // here the value of i was passed into as the argument cntr
        // and will be captured in this function closure so each
        // iteration of the loop can have it's own value
    })(i);
}
