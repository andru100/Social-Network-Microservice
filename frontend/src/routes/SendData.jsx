export default async function SendDdata (Page, msgType, reply, Liker, iD, Post) { // send comment or reply to comment depending on is replychk false/true

    let data = {Username: Page, Updatetype: '$set',
        Key2updt: 'Posts', isReply: msgType, replyCmt: reply, LikeSent: Liker, PostIndx: iD, Value2updt: Post
    }

    let options = {
    method: 'PUT',
    headers: {
    'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
    }

    let postUrl = process.env.REACT_APP_BACKEND_ADDRESS + process.env.REACT_APP_POSTMSG_PORT + '/postMsg'
    let response = await fetch(postUrl, options)
    let convert = await response.json ()

    if ( response.status === 401 || response.status === 400){
        console.log("error sending comment")
    } else if ( response.status === 200){ 
        return convert
    }

}