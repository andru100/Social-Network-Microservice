import SendData from "./SendData"

export default async function RequestOTP (username, requesttype, usertype) {

    console.log("Requesting OTP", requesttype)
        
    let signindata = {data: {
        Username: username,
        RequestType: requesttype,
        UpdateType: usertype
        }
    }

    let gqlRequest = "mutation SecureUpdate ($data: SecurityCheckInput!){ SecureUpdate(input: $data) { Token AuthType MobClue EmailClue } }"
   // let gqlRequest = "mutation SignIn ($data: UsrsigninInput!){ SignIn(input: $data) { Token } }"
    
    let response = await SendData(gqlRequest, signindata, 'secureupdate')

    return response

}

