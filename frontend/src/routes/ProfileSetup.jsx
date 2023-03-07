import { useEffect, useState} from 'react'
import { useNavigate, useParams } from 'react-router-dom';
import ChkAuth from '../chkAuth';
import bkpic from '../images/profbkup.jpg';
  
export default function RenderProfileSetup () { 
  const [sessionUser, setSessionUser] = useState ("")
  const[cmt, setcmt] = useState({Posts:[], Photos: []})
  const [dp, setDp] = useState(false); // hide/show div
  const Navigate = useNavigate();
  var {user} = useParams()
  let userName = user

  useEffect( () => {
    ChkAuth().then(user => {
      if (user) {
        setSessionUser(user)
        getCmt().then(cmtz => {
          if (cmtz) {
          setcmt(cmtz)
          console.log("returned users data object", cmtz)
          }
        })
      } else {
         setSessionUser(false)
         alert("You need to sign in to view this page")
         Navigate("/signIn")
      }
    })
  },[]);

  function goToProfile (){
    Navigate ("/profile/x/"+sessionUser)
  }
  
  function Logout (){
    localStorage.removeItem('jwt_token')
    Navigate ("/")
  }

  async function getCmt () { // gets users data object
  
    let data = {Username: userName, Data: ""}

    let options = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
    }

    let Url = process.env.REACT_APP_BACKEND_ADDRESS + process.env.REACT_APP_STDGET_PORT + '/getCmt'
    let response = await fetch(Url, options)
    let convert = await response.json ()

    if ( response.status === 401 || response.status === 400){// handle not a match 
      console.log("error fetching, returned data is:", response)
    } else if ( response.status === 200){
      console.log("found users data:", convert)
      return convert
    }
   
  }


  async function updateBio () { 
    
    const bio = document.getElementById('bioBox').value;// get username, password and email.

    let data = {Username: userName, Updatetype: '$set',
                Key2updt: 'Bio', Value2updt: bio
    }

    let options = {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
    }

    let postUrl = process.env.REACT_APP_BACKEND_ADDRESS + process.env.REACT_APP_UPDATEBIO_PORT + '/updatebio'
    let response = await fetch(postUrl, options)
    let convert = await response.json ()

    if ( response.status === 401 || response.status === 400){
      alert("your bio didnt save try again")
    } else if ( response.status === 200){
      console.log("saved bio")
      setcmt(convert) // store users data object
      setDp(!dp) // show bio edit box
    }
  }

  async function addPhotos (event) {
    if (event) {
      alert("in add photos")
      let file = (event.target.files[0])
      var data = new FormData() 
      data.append('file', file)
      data.append('user', userName)
      data.append('type', 'addPhotos')
    
      let options = {
        method: 'POST',
        body: data, 
      }

      let postUrl = process.env.REACT_APP_BACKEND_ADDRESS + process.env.REACT_APP_POSTFILE_PORT + '/postfile/' + userName 
      let response = await fetch(postUrl, options)
      let convert = await response.json ()
      if ( response.status === 401 || response.status === 400){
        alert("your pic didn't save try again")
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
      data.append('user', userName)
      data.append('type', 'profPic')
    
      let options = {
        method: 'POST',
        body: data, 
      }

      let ProfUrl = process.env.REACT_APP_BACKEND_ADDRESS + process.env.REACT_APP_POSTFILE_PORT + '/postfile/' + userName
      let response = await fetch(ProfUrl, options)
      let convert = await response.json ()
      document.getElementById("profpic11").src = convert.Profpic // get posted img address and change profile picture
      
    }
  };

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
              <a href="!" className="btn btn-sm btn-info mb-2" style={{marginRight:"10px"}} onClick={(e) => {e.preventDefault(); goToProfile()}}>View Profile</a>
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
              {sessionUser === user &&
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