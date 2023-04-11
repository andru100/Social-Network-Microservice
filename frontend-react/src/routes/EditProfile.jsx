import { useEffect, useState} from 'react'
import { useNavigate, useParams } from 'react-router-dom';
import ChkAuth from './chkAuth';
import SendData from './SendData';
import bkpic from '../images/profbkup.jpg';
import Home from './Home';
  
export default function RenderProfileSetup (props) { 
  const [sessionuser, setSessionUser] = useState (props.sessionuser)
  const[cmt, setcmt] = useState({Posts:[], Photos: []})
  const [dp, setDp] = useState(false); // hide/show div
  const [page, setPage] = useState("default"); // hide/show div
  let Page =  "Bio"

  console.log("in RenderProfileSetup, sessionuser:", sessionuser)

  useEffect( () => {
        getCmt().then(cmtz => {
          if (cmtz) {
          setcmt(cmtz)
          console.log("use effet get coment users data object is :", cmtz)
        }   
    })
  },[]);

  function GoHome () {
		console.log("going home, username:", sessionuser)
		setPage("home")
	}
  
  function Logout (){
    localStorage.removeItem('jwt_token')
    //Navigate ("/")
  }

  async function getCmt (user2find) { // sends Username, password from input, then backend creates s3 bucket in Username and stores details on mongo
  
    let data = {Username: sessionuser}

    let queryType

    if (Page === "All") {
       queryType = "GetAllComments"
    } else {
       queryType ="GetUserComments"
    }


    if (user2find) {
       data.Username = sessionuser
    }


    let gqlRequest = "query " + queryType + " ($Username: String!){  " + queryType + " (input: $Username) { Key ID Username Password Email Bio Profpic Photos LastCommentNum Posts { Username SessionUser MainCmt PostNum Time TimeStamp Date Comments { Username Comment Profpic } Likes { Username Profpic } } } }"
    let response = await SendData(gqlRequest, data)
    if ( "errors" in response ){ // if password is a match redirect to profile page
      //{ProcessErrorAlerts("hi", "hi")}
      console.log("error retrieving data", response.errors[0].message )
      return false
      
    } else {
       return response.data[queryType] 
    }
 }


  async function updateBio () { 
    
    const bio = document.getElementById('bioBox').value

    let variables = {data: {Username: sessionuser, Bio: bio}}

    let gqlRequest = "mutation UpdateBio ($data: UpdateBioInput!){ UpdateBio(input: $data) { Key ID Username Password Email Bio Profpic Photos LastCommentNum Posts { Username SessionUser MainCmt PostNum Time TimeStamp Date Comments { Username Comment Profpic } Likes { Username Profpic } } } }"
		
		let response = await SendData(gqlRequest, variables)

		
		if ( "errors" in response ){ // if password is a match redirect to profile page
      //{ProcessErrorAlerts("hi", "hi")}
      console.log("error updating bio", response.errors[0].message )
      return false
      
    } else { // if password is a match redirect to profile page
			console.log("saved bio")
      setcmt(response.data.UpdateBio) // store users data object
      setDp(!dp) // show bio edit box
		} 
  }

  async function addPhotos (event) {
    if (event) {
      let file = (event.target.files[0])
      var data = new FormData() 
      data.append('file', file)
      data.append('user', sessionuser)
      data.append('type', 'addPhotos')
    
      let options = {
        method: 'POST',
        body: data, 
      }

      let postUrl = process.env.REACT_APP_BACKEND_UPLOAD +  '/postfile/' + sessionuser 
      let response = await fetch(postUrl, options)
      let convert = await response.json ()
      if ( response.status === 401 || response.status === 400){
        console.log("your pic didn't save, please try again")
       } else if ( response.status === 200){ 
        console.log("added pic to users photos")
        setcmt(convert)
       }
      
    }
  };

  function triggerClick(event){ // clicking image triggers upload button click
    var myButton = document.getElementById(event.target.name);
    if ( myButton ) {
      myButton.click()
    }
  }

  async function addProfilePic (event) {
    if (event) {
      let file = (event.target.files[0])
      var data = new FormData() 
      data.append('file', file)
      data.append('user', sessionuser)
      data.append('type', 'profPic')
    
      let options = {
        method: 'POST',
        body: data, 
      }

      let ProfUrl = process.env.REACT_APP_BACKEND_UPLOAD + '/postfile/' + sessionuser
      let response = await fetch(ProfUrl, options)
      let convert = await response.json ()
      document.getElementById("profpic11").src = convert.Profpic // get posted img address and change profile picture
      
    }
  };

  function LandingPage() {
    return (
      <>
        <span className="login100-form-title p-b-53">
          Create your profile
        </span>
          {/*  <!--profile pic--> */}
          <div className="row">
            <div class="col"></div>
            <div class="col"></div>
            <div class="col"></div>
            <div class="col"></div>
              <div className="col" >
                <div className="profpics">
                  {cmt.Profpic ? <img className="profpics" id= "profpic11" name="profpic1"  onClick={(e)=>  {e.preventDefault(); triggerClick(e)}} alt="alt" src={cmt.Profpic}  
                    data-holder-Rendered="true"/> 
                    :
                    <img className="profpics" id= "profpic11" name="profpic1"  onClick={(e)=>  {e.preventDefault(); triggerClick(e)}} alt="alt" src={bkpic}  
                    data-holder-Rendered="true"/>
                  }
                </div>
              </div>
            <div className="visually-hidden">
                <div>
                      <input id="profpic1" type="file" className="blocked" onChange={addProfilePic}  name= "uploader1"/>
                </div>
            </div>
              <div className="col-md-4 mb-3" >
            </div>
            <div className="visually-hidden"></div>
          </div>    
        <div className="w-full text-center p-t-55">
        <a href="!" className="btn btn-sm btn-info mb-2" style={{marginRight:"10px"}} onClick={(e) => {e.preventDefault(); GoHome()}}>View Profile</a>
        <a href="!" className="btn btn-sm btn-info mb-2" onClick={(e) => {e.preventDefault(); Logout()}}>Log Out</a>
        </div>               
        <div className="p-t-13 p-b-9"></div>
        <div className="wrap-input100 validate-input">{cmt.Bio ? cmt.Bio : "You haven't added a bio yet. Click to add one now!"}</div>
        <a href="!" className="btn btn-sm btn-info mb-2" style={{marginTop: "10px"}} onClick={(e) => {e.preventDefault(); setDp(!dp)}}>Edit Bio</a>
        {dp && 
        <>
          <div className="wrap-input100 validate-input" style={{ display: dp }}>
            ​<span className="txt10">
              Max 70 characters
            </span> 
            <span><textarea maxLength="70" id="bioBox" rows="10" cols="10"></textarea></span>
          </div>
          <div className="container-login100-form-btn m-t-17">
            <button className="login100-form-btn" type="button" onClick={(e) =>{e.preventDefault(); updateBio()}}>
              Save Changes
            </button>
          </div>
        </>
        }
        {sessionuser === sessionuser &&
        <div className="w-full text-center p-t-55">
          <div className="connected-container">
            <div className="gif-grid">
            {/* map through users images */}
                {cmt.Photos.map((pic) => (
                  <div className="gif-item" key={pic}>
                      <img src={pic} alt={pic} />
                  </div>
                ))}
            </div>
          </div>
        </div>
        }
        <div>
          <a href="!" className="btn btn-sm btn-info mb-2" name="mediaUplaod" style={{marginTop: "10px"}} onClick={(e) => {e.preventDefault(); triggerClick(e)}}>Add Photos</a>
        </div>  
        <div>
            <input id="mediaUplaod" type="file" className="blocked" onChange={(e)=> addPhotos(e)}  name= "uploader1"/>
        </div>
        <div>
        </div>
      </>
    )
  }

  return ( 
    <>
      
      {page === "default" && <LandingPage />}
      {page === "home" &&  <Home sessionuser={sessionuser} viewing={sessionuser} page={"home"}/>}
      
    </>
    )

}