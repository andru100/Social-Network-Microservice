export default async function SendData (request, variables ) { // send data to GraphQL
    
    let GQLpayload = {
        "query": request,
        "variables":variables
    }
   
    let options = {
        method: 'POST',
        headers: {
        'Content-Type': 'application/json',
        },
        body: JSON.stringify(GQLpayload)
    }

    let postUrl = process.env.REACT_APP_BACKEND_GRAPHQL + "/query" 
    let response = await fetch(postUrl, options)
    let convert = await response.json ()

    console.log("response from GraphQL", convert)
    
    
    return convert
    

}