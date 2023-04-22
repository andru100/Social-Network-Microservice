import { React, memo } from "react";
import { useState, useEffect } from "react"
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import ChkAuth from './chkAuth';
import SendData from './SendData';
import UpdateHybrid from "./UpdateHybrid";
import SignIn from './SignIn'

function Home (props) {
   const [sessionUser, setSessionUser] = useState (props.sessionuser)
   const[cmt, setcmt] = useState({Posts:[], Photos: []}) // holds users data object
   const [viewReply, setviewReply] = useState({}); // use to show comments when clicked
   const [viewLikes, setviewLikes] = useState({}); // use to show likes when clicked
   const [viewCmtBox, setviewCmtBox] = useState({}); // use to show comment box when clicked
   const [page, setPage] = useState(props.page); // use to show comment box when clicked
   const [viewing, setViewing] = useState(props.viewing); // use to show comment box when clicked
   const [scope, setScope] = useState("user"); // use to show comment box when clicked

   console.log("render occcured, scope is: ", scope, "viewing: ", viewing, "sessionUser: ", sessionUser, "page: ", page)
   console.log("user data is: ", cmt)

   dayjs().format()
   dayjs.extend(relativeTime)

   var timeAtRender = dayjs(Date.now())

  useEffect( () => {
      getCmt("user").then(cmtz => {
         if (cmtz) {
         setcmt(cmtz)
         console.log("Users data object retrieved is:", cmtz)
         }
      })
  },[]);

// extra auth heck option
//   useEffect( () => {
//    ChkAuth().then(user => {
//          if (user) {
//             setSessionUser(user)
//             getCmt().then(cmtz => {
//                if (cmtz) {
//                setcmt(cmtz)
//                console.log("Users data object retrieved is:", cmtz)
//                }
//             })
//          } else {
//             setSessionUser(false)
//             alert.show("You need to sign in to view this page")
//             setPage("signin")
//          }
//       })
//   },[]);


   async function getCmt (request, user, searchKey) { // sends username, password from input, then backend creates s3 bucket in username and stores details on mongo
  
      let payload = { data: {
         Username: viewing,
         RequestType: request
      }}

      user && (payload.data.Username = user)
      searchKey && (payload.data.SearchTerm = searchKey)


      let gqlRequest = "query GetPosts ($data: GetPost!){  GetPosts (input: $data){ Username Bio Profpic Photos Following Followers Replys{Username PostID ReplyID} Liked{Username PostID} Posts { ID Username MainCmt TimeStamp Comments {ID Username Comment Profpic } Likes {Username Profpic } } } }"
      let response = await SendData(gqlRequest, payload)
      if ( "errors" in response ){ // if password is a match redirect to profile page
			//{ProcessErrorAlerts("hi", "hi")}
			console.log("Error retrieving user data", response.errors[0].message )
			
		} else {
         return response.data["GetPosts"] 
      }
   }

   async function getSessionUserData () { // sends username, password from input, then backend creates s3 bucket in username and stores details on mongo
  
      let payload = {data: {
         Username: sessionUser,
         RequestType: "user"
      }}
      


      let gqlRequest = "query GetPosts ($data: GetPost!){ GetPosts (input: $data){ Username Bio Profpic Photos Following Followers Replys{Username PostID ReplyID} Liked{Username PostID} Posts { ID Username MainCmt TimeStamp Comments {ID Username Comment Profpic } Likes {Username Profpic } } } }"
      let response = await SendData(gqlRequest, payload)
      if ( "errors" in response ){ 
			console.log("Error retrieving user data", response.errors[0].message )
		} else {
         return response.data["GetPosts"] 
      }
   }


   async function sendCmt (cmtAuthr) { // sends comments, replies to comments and likes
    
      let NewCmtInput  = { data: {
         Username : cmtAuthr,
         SessionUser: sessionUser, //  not needed once can be retrieved from the session cookie
         MainCmt : document.getElementById("cmt").value ,
         TimeStamp : Date.now(),    
         ReturnPage: "user"   ,
         RequestType: "create"
         }
      }

      let gqlRequest = "mutation NewComment ($data: SendCmtInput!){ NewComment (input: $data){ Username Bio Profpic Photos Following Followers Replys{Username PostID ReplyID} Liked{Username PostID} Posts { ID Username MainCmt TimeStamp Comments {ID Username Comment Profpic } Likes {Username Profpic } } } }"
      
      SendData(gqlRequest, NewCmtInput).then((response)=> ("errors" in response) ? console.log("Error posting comment", response.errors[0].message ) : setcmt(response.data.NewComment) ); 
   
   }

   async function sendReply (cmtAuthr, iD) { // sends comments, replies to comments and likes
      
      let CommentResponse =  { data: {
         AuthorUsername: cmtAuthr ,
         ReplyUsername: sessionUser ,
         ReplyComment:  document.getElementById(iD).value , 
         ReplyProfpic:  "" ,
         PostID:   iD ,
         ReturnPage: scope,
         RequestType: "create"
         }
      }

      let gqlRequest = "mutation ReplyComment ($data: ReplyCommentInput!){ ReplyComment (input: $data){ Username Bio Profpic Photos Following Followers Replys{Username PostID ReplyID} Liked{Username PostID} Posts { ID Username MainCmt TimeStamp Comments {ID Username Comment Profpic } Likes {Username Profpic } } } }"
      const reply = await getSessionUserData()
      CommentResponse.data.ReplyProfpic = reply.Profpic
      
      SendData(gqlRequest, CommentResponse).then((response)=> ("errors" in response) ? console.log("error sending response to comment") : setcmt(response.data.ReplyComment))
   }

   async function sendLike (cmtAuthr, iD) { 

      let SendLikeInput  = { data: {
         Username:   cmtAuthr ,
         LikedBy:   sessionUser ,
         LikeByPic:   "",
         PostID:   iD , 
         ReturnPage: scope , 
         RequestType: "create"
         }
      }
      let gqlRequest = "mutation LikeComment ($data: SendLikeInput!){ LikeComment (input: $data){ Username Bio Profpic Photos Following Followers Replys{Username PostID ReplyID} Liked{Username PostID} Posts { ID Username MainCmt TimeStamp Comments {ID Username Comment Profpic } Likes {Username Profpic } } } }"
      getSessionUserData().then((repliersData)=> {SendLikeInput.data.LikeByPic = repliersData.Profpic; SendData(gqlRequest, SendLikeInput).then((response)=> ( "errors" in response) ? console.log("error adding like") : setcmt(response.data.LikeComment) ) ; })
   }

   async function deleteCmt (cmtAuthr, postID) { // sends comments, replies to comments and likes

      let NewCmtInput  = { data: {
         Username : cmtAuthr,
         SessionUser: sessionUser,  
         ReturnPage: scope   ,
         RequestType: "delete",
         PostID: postID
         }
      }

      let gqlRequest = "mutation NewComment ($data: SendCmtInput!){ NewComment (input: $data){ Username Bio Profpic Photos Following Followers Replys{Username PostID ReplyID} Liked{Username PostID} Posts { ID Username MainCmt TimeStamp Comments {ID Username Comment Profpic } Likes {Username Profpic } } } }"
      
      SendData(gqlRequest, NewCmtInput).then((response)=> ("errors" in response) ? console.log("Error posting comment", response.errors[0].message ) : setcmt(response.data.NewComment) ); 
   
   }

   async function deleteResponse (cmtAuthr, postID, replyID) { // sends comments, replies to comments and likes
      
      let CommentResponse =  { data: {
         AuthorUsername: cmtAuthr ,
         ReplyUsername: sessionUser ,
         PostID:   postID ,
         ReplyID: replyID, 
         ReturnPage: scope ,
         RequestType: "delete",
         }
      }

      let gqlRequest = "mutation ReplyComment ($data: ReplyCommentInput!){ ReplyComment (input: $data){ Username Bio Profpic Photos Following Followers Replys{Username PostID ReplyID} Liked{Username PostID} Posts { ID Username MainCmt TimeStamp Comments {ID Username Comment Profpic } Likes {Username Profpic } } } }"
      
      SendData(gqlRequest, CommentResponse).then((response)=> ("errors" in response) ? console.log("error sending response to comment") : setcmt(response.data.ReplyComment))
   
   }

   async function deleteLike (cmtAuthr, iD) { // sends comments, replies to comments and likes

      let SendLikeInput  = { data: {
         Username:   cmtAuthr ,
         LikedBy:   sessionUser ,
         PostID:   iD , 
         ReturnPage: scope ,  
         RequestType: "delete"
         }
      }
      let gqlRequest = "mutation LikeComment ($data: SendLikeInput!){ LikeComment (input: $data){ Username Bio Profpic Photos Following Followers Replys{Username PostID ReplyID} Liked{Username PostID} Posts { ID Username MainCmt TimeStamp Comments {ID Username Comment Profpic } Likes {Username Profpic } } } }"
      SendData(gqlRequest, SendLikeInput).then((response)=> ( "errors" in response) ? console.log("error adding like") : setcmt(response.data.LikeComment) ) 
   

   }

   async function follow (requestype) { // sends comments, replies to comments and likes

      let SendFollowInput  = { data: {
         Username:  sessionUser ,
	      UserOfIntrest: viewing ,
	      ReturnPage: "user" ,
	      RequestType: requestype
         }
      }
      let gqlRequest = "mutation Follow ($data: FollowInput!){ Follow (input: $data) { Username Bio Profpic Photos Following Followers Replys{Username PostID ReplyID} Liked{Username PostID} Posts { ID Username MainCmt TimeStamp Comments {ID Username Comment Profpic } Likes {Username Profpic } } } }"
      SendData(gqlRequest, SendFollowInput).then((response)=> ( "errors" in response) ? console.log("error adding like") : setcmt(response.data.Follow) ) 
   }




   //edit profile functions

   async function updateBio () { 
    
      const bio = document.getElementById('bioBox').value
  
      let variables = {data: {Username: sessionUser, Bio: bio}}
  
      let gqlRequest = "mutation UpdateBio ($data: UpdateBioInput!){ UpdateBio(input: $data) { Username Bio Profpic Photos Posts { ID Username MainCmt TimeStamp Comments {ID Username Comment Profpic } Likes {Username Profpic } } } }"
        
        let response = await SendData(gqlRequest, variables)
  
        
        if ( "errors" in response ){ // if password is a match redirect to profile page
        //{ProcessErrorAlerts("hi", "hi")}
        console.log("error updating bio", response.errors[0].message )
        return false
        
      } else { // if password is a match redirect to profile page
           console.log("saved bio")
        setcmt(response.data.UpdateBio) // store users data object
        //setDp(!dp) // show bio edit box
        //change to update details type or put in unpdate details
        } 
    }
  
    async function addPhotos (event) {
      if (event) {
        let file = (event.target.files[0])
        var data = new FormData() 
        data.append('file', file)
        data.append('user', sessionUser)
        data.append('type', 'addPhotos')
      
        let options = {
          method: 'POST',
          body: data, 
        }
  
        let postUrl = process.env.REACT_APP_BACKEND_UPLOAD +  '/postfile/' + sessionUser 
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
        data.append('user', sessionUser)
        data.append('type', 'profPic')
      
        let options = {
          method: 'POST',
          body: data, 
        }
  
        let ProfUrl = process.env.REACT_APP_BACKEND_UPLOAD + '/postfile/' + sessionUser
        let response = await fetch(ProfUrl, options)
        let convert = await response.json ()
        document.getElementById("profpic11").src = convert.Profpic // get posted img address and change profile picture
        
      }
    };


    function goToHome (){
      getCmt("user", sessionUser).then(cmtz => {
         if (cmtz) {
          setcmt(cmtz)
          console.log("Users data object retrieved is:", cmtz)
         }
      })
      setScope("user")
      setViewing(sessionUser)
    }

    function goToUser (user){
      getCmt("user", user).then(cmtz => {
         if (cmtz) {
          setcmt(cmtz)
          console.log("Users data object retrieved is:", cmtz)
         }
      })
      setScope("user")
      setViewing(user)
      
    }

    function goToSuggestedPosts (){
      getCmt("suggested", sessionUser).then(cmtz => {
         if (cmtz) {
          setcmt(cmtz)
          console.log("Users data object retrieved is:", cmtz)
         }
      })
      setScope("suggested")
      // getPosts()
    }

    function goToFriends (){
      getCmt("following", sessionUser).then(cmtz => {
         if (cmtz) {
          setcmt(cmtz)
          console.log("get friends object retrieved is:", cmtz)
         }
      })
      setScope("following")
    }

    function goToReplys (){
      getCmt("replys", viewing).then(cmtz => {
         if (cmtz) {
          setcmt(cmtz)
          console.log("get friends object retrieved is:", cmtz)
         }
      })
      //setScope("user")
    }

    function goToLikes (){
      getCmt("likes", viewing).then(cmtz => {
         if (cmtz) {
          setcmt(cmtz)
          console.log("get friends object retrieved is:", cmtz)
         }
      })
      //setScope("user")
    }

    function goToPhotos () {
      //  if (page !== "media") { // for when your on media tab already
         setScope("media")
   
    }

    function Logout (){
      localStorage.removeItem('jwt_token')
      setPage("signin")
      //Navigate ("/")
    }

    const toggleLikes = (id) => {
      setviewLikes(prev => Boolean(!prev[id]) ? {...prev, [id]: true} : {...prev, [id]: false});
    };

    const toggleReply = (id) => {
      setviewReply(prev => Boolean(!prev[id]) ? {...prev, [id]: true} : {...prev, [id]: false});
    };

    const toggleCmt = (id) => {
      setviewCmtBox(prev => Boolean(!prev[id]) ? {...prev, [id]: true} : {...prev, [id]: false});
    };

    const containerStyle= {
      width: '100%',  
      minHeight: '100vh',
      backgroundPosition: 'center',
      backgroundSize: 'cover',
      backgroundRepeat: 'no-repeat',
      backgroundImage: 'linear-gradient(47deg, cyan, magenta)'
   }

   function Home() {
      return (
        <div style={containerStyle}>
         <Container>
           <Row>
               <Col md="3"></Col>
               <Col md="6">
                  <div className="comments" style={{background: "black"}}>
                     <Row>
                        <Col style={{marginLeft: "30px"}}>
                           {/* <div className="profile-header-content"> */}
                              <div className="profile-header-img">
                                 {/* <img className="profpics" id= "profpic11" name="profpic1"  onClick={(e)=>  {e.preventDefault(); triggerClick(e)}} alt="alt" src={cmt.Profpic}  data-holder-Rendered="true"/> */}
                                 <img calssName="profpics" id= "profpic11" name="profpic1" src={cmt.Profpic} alt="alt" data-holder-Rendered="true" onClick={(e)=>  {e.preventDefault(); triggerClick(e)}}/>
                              </div>
                              <div className="visually-hidden">
                                 <div>
                                       <input id="profpic1" type="file" className="blocked" onChange={addProfilePic}  name= "uploader1"/>
                                 </div>
                              </div>
                           {/* </div> */}
                        </Col>
                        <Col>
                           <div className="profile-header-info">
                              <h4 className="m-t-10 m-b-5">{viewing}</h4>
                              <p className="m-b-10" style={{color:"black"}}>{cmt.Bio? cmt.Bio : "Click the edit profile button to add a bio to your profile now."}</p>
                              {
                                 sessionUser === viewing ?

                                 <a href="!" className="btn btn-sm btn-info mb-2" style={{marginRight:"10px"}} onClick={(e)=> {e.preventDefault() ; setPage("updatedetails")}}>Update Details</a>

                                 :

                     
                                 
                                 cmt.Followers && cmt.Followers.includes(sessionUser) ?
                                 
                                 <a href="!" className="btn btn-sm btn-info mb-2" style={{marginRight:"10px"}} onClick={(e)=> {e.preventDefault() ; follow("unfollow")}}>Unfollow</a>


                                 :
                                 <a href="!" className="btn btn-sm btn-info mb-2" style={{marginRight:"10px"}} onClick={(e)=> {e.preventDefault() ; follow("follow")}}>Follow</a>
                                
                              }
                              
                              
                              <a href="!" className="btn btn-sm btn-info mb-2" onClick={(e)=> {e.preventDefault(); Logout()}}>Log Out</a>
                           </div>
                        </Col>
                     </Row>
                  </div>
               </Col>
               <Col md="3"></Col>
                              
            </Row>


             <Row>
               <Col md="3"></Col>
               <Col md="6">
                  <Row>
                     <Col>
                        {scope === "user" && <button className="login100-form-btn" type="button" onClick={(e)=> {e.preventDefault() ; goToUser(viewing)}}>Posts</button>}
                     </Col>
                     <Col>
                        {scope === "user" && <button className="login100-form-btn" type="button" onClick={(e)=> {e.preventDefault() ; goToReplys(viewing)}}>Replys</button>}
                     </Col>
                     <Col>
                        {scope === "user" && <button className="login100-form-btn" type="button" onClick={(e)=> {e.preventDefault() ; goToLikes(viewing)}}>Likes</button>}
                     </Col>
                     <Col>
                        {scope === "user" && <button className="login100-form-btn" type="button" onClick={(e)=> {e.preventDefault() ; goToPhotos(viewing)}}>Media</button>}
                     </Col>
                  </Row>
               </Col>
               <Col md="3"></Col>
                              
            </Row>      
              
             
           <Row>
             <Col xs lg="3"> 
               <Row>
                   <Col style={{}}>
                     {sessionUser && <button className="login100-form-btn" type="button" onClick={(e)=> {e.preventDefault() ; goToHome()}}>PROFILE</button>}
                   </Col>
               </Row>
               <Row>
                     <Col style={{}}>
                     {sessionUser && <button className="login100-form-btn" type="button"  onClick={(e)=> {e.preventDefault() ; goToSuggestedPosts()}}>FOR YOU</button>}
                             
                     </Col>
               </Row>
               <Row>
                     <Col style={{}}>
                     {sessionUser && <button className="login100-form-btn" type="button"  onClick={(e)=> {e.preventDefault() ; goToPhotos()}}>MEDIA</button>}      
                     </Col>
               </Row>
               <Row>
                     <Col>
                     {sessionUser && <button className="login100-form-btn" type="button" onClick={(e)=> {e.preventDefault() ; goToFriends()}}>FRIENDS</button>}  
                     </Col>        
               </Row>
             </Col>
             <Col>
                  {sessionUser && scope === "media" &&
                  <>
                     <div className="connected-container">
                        <div className="gif-grid">
                           {cmt.Photos.map((pic) => (
                              <div className="gif-item" key={pic}>
                                 <img src={pic} alt={pic} />
                              </div>
                           ))}
                        </div>
                     </div>
                     <div>
                        <a href="!" className="btn btn-sm btn-info mb-2" name="mediaUplaod" style={{marginTop: "10px"}} onClick={(e) => {e.preventDefault(); triggerClick(e)}}>Add Photos</a>
                     </div>  
                     <div>
                           <input id="mediaUplaod" type="file" className="blocked" onChange={(e)=> addPhotos(e)}  name= "uploader1"/>
                     </div>
                  </>
                  }

                  
                   {/*renders the post comment box if user is logged in and viewing theyre own profile or all posts page*/}
                   {(sessionUser === viewing && page !== "media") || (sessionUser && page === "all") ? 
                       <div className= "comments">
                        <Row>
                           
                              <form action="">
                                 <div className="input-group">
                                    <textarea style={{height: "36px"}} type="text" className="form-control rounded-corner" id="cmt" placeholder="Write a comment..."/>
                                    <span className="input-group-btn p-l-10">
                                       <button className="btn btn-primary f-s-12 rounded-corner" type="button" onClick={() => sendCmt(sessionUser)}>Comment</button>
                                    </span>
                                 </div>
                              </form>
                           
                        </Row>
                       </div> : null}
                        


                  {/*render users profile page or news feed showing all comments*/}
                  { page !== "media" &&
                     cmt.Posts.map((userData)=> {
                        return (
                           <>
                              <div className="comments" style={{background: "white"}}>
                                 <Row>
                                   <div>
                                    <Row>
                                       <Col md="auto">
                                          <span className="userimage"><img onClick={()=>   goToUser(userData.Username)} src={cmt.Profpic} alt=""/></span>
                                       </Col>
                                       <Col>
                                          <Row>
                                             <Col md="auto">
                                                <div className="username" onClick={()=> goToUser(userData.Username)}>{userData.Username}</div>
                                             </Col>
                                             <Col>
                                                <div className="time">{dayjs(userData.TimeStamp).from(timeAtRender) }</div>
                                             </Col>
                                             <Col>
                                                {
                                                   sessionUser === userData.Username &&

                                                   <i class="fa fa-trash-o" onClick={() => deleteCmt(sessionUser, userData.ID)} aria-hidden="true"></i>
                                                }
                                             </Col>

                                                
                                                
                                          </Row>
                                          <Row>
                                             <Col>
                                                <div className="comment-content">
                                                   <p>
                                                      {userData.MainCmt}
                                                   </p>
                                                </div>
                                             </Col>
                                          </Row>
                                       </Col>
                                    </Row>
                                    <Row>
                                       <Col>
                                          <div className="reply-icons" >
                                             { 
                                                userData.Likes.some(e => e.Username === sessionUser) ?
                                                   <span className="fa-stack fa-fw stats-icon">
                                                      <i className="fa fa-circle fa-stack-2x text-danger"></i>
                                                      <i className="fa fa-heart fa-stack-1x fa-inverse t-plus-1" onClick={()=>deleteLike(userData.Username, userData.ID)}></i>
                                                   </span>

                                                :

                                                <span className="fa-stack fa-fw stats-icon">
                                                   <i className="fa fa-heart fa-stack-1x fa-inverse t-plus-1" style={{backgroundColor: "red"}} onClick={()=>sendLike(userData.Username, userData.ID)}></i>
                                                </span>

                                             }
                                             <span className="stats-text" onClick={() => {viewReply[userData.PostNum] && toggleReply(userData.PostNum) ; viewCmtBox[userData.PostNum] && toggleCmt(userData.PostNum); toggleLikes(userData.PostNum)}}>{userData.Likes?.length} Likes</span>
                                          </div>
                                       </Col>
                                       <Col>
                                             <i class="fa fa-comment-o" aria-hidden="true"></i>
                                             <span className="stats-text" onClick={() => {viewLikes[userData.PostNum] && toggleLikes(userData.PostNum) ; viewCmtBox[userData.PostNum] && toggleCmt(userData.PostNum); toggleReply(userData.PostNum)}}>{userData.Comments?.length} Comments</span>
                                       </Col>
                                       <Col>
                                          <i class="fa fa-reply" aria-hidden="true"></i>
                                          <span className="stats-text" onClick={() => { viewLikes[userData.PostNum] && toggleLikes(userData.PostNum) ; viewReply[userData.PostNum] && toggleReply(userData.PostNum) ;  toggleCmt(userData.PostNum)}}>Reply</span>
                                       
                                       </Col>
                                    </Row>
                                    <Row>
                                             {/* show likes */}
                                             {viewLikes[userData.PostNum] &&
                                             userData.Likes.map((Likes)=> (
                                                   <Row >
                                                      <div className="replys">
                                                         <Col md="auto">
                                                            <span className="userimage"><img onClick={()=> setViewing(userData.Username)} src={Likes.Profpic} alt=""/></span>  
                                                         </Col>
                                                         <Col>
                                                            <Row>
                                                               <div ClassName = "username" onClick={()=> setViewing(userData.Username)} >{Likes.Username}</div> 
                                                            </Row>
                                                         </Col>
                                                      </div>
                                                      
                                                   </Row> 
                                             ))
                                             } 

                                             {/* show comments */}
                                             {viewReply[userData.PostNum] &&
                                             userData.Comments.map((replys)=> (
                                                <Row >
                                                   <div className="replys">
                                                      <Col  md="auto">
                                                         <span className="userimage"><img onClick={()=> setViewing(userData.Username)} src={replys.Profpic} alt=""/></span>  
                                                      </Col>
                                                      <Col>
                                                         <Row>
                                                            <Col>
                                                                <div ClassName = "username" onClick={()=> setViewing(replys.Username)} >{replys.Username}</div>
                                                            </Col> 
                                                            <Col>
                                                            {
                                                               sessionUser === replys.Username &&

                                                               <i class="fa fa-trash-o" onClick={() => deleteResponse(userData.Username, userData.ID, replys.ID)} aria-hidden="true"></i>
                                                            }
                                                            </Col>
                                                         </Row>
                                                         <Row>
                                                            <p>{replys.Comment}</p>
                                                         </Row>
                                                      </Col>
                                                   </div>
                                                </Row>
                                             ))
                                             } 

                                             {/* show reply box */}
                                             {viewCmtBox[userData.PostNum] && 
                                                <>
                                                   <Row>
                                                         <div className="input">
                                                            <form action="">
                                                               <div className="input-group">
                                                                  <textarea type="text" className="form-control rounded-corner" id={userData.ID} style={{height:"90px", width:"180px"}} placeholder="Reply to the post..."/>
                                                               </div>
                                                            </form>
                                                         </div>
                                                   </Row>
                                                   <Row>
                                                      <span className="input-group-btn p-l-10">
                                                         <button className="btn btn-primary f-s-12 rounded-corner" type="button"  onClick={()=>sendReply(userData.Username, userData.ID)}>Comment</button>
                                                      </span>
                                                   </Row>
                                                </> 
                                             }
                                            
                                    </Row>
                                  </div>
                                    
                                 </Row>
                              </div>
                           </>
                  )})}
                              

                  
             </Col>
             <Col xs lg="3">
             </Col>
          
           </Row>
         </Container>
      </div>
       );
   }

   return (
      <>
      {page === "updatedetails" && <UpdateHybrid username={sessionUser} /* setSessionUser={setSessionUser()} */ />}
	   {page === "home" &&  <Home sessionuser={sessionUser} page={page} viewing={sessionUser} />}
      {page === "signin" && <SignIn/>}
      </>

   )
}

export default memo(Home)