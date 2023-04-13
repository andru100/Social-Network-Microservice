import { React, memo } from "react";
import { useState, useEffect } from "react"
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import ChkAuth from './chkAuth';
import SendData from './SendData';
import UpdateDetails from './UpdateDetails'
import EditProfile from './EditProfile'
import SignIn from './SignIn'

function Home (props) {
   const [sessionUser, setSessionUser] = useState (props.sessionuser)
   const[cmt, setcmt] = useState({Posts:[], Photos: []}) // holds users data object
   const [viewReply, setviewReply] = useState({}); // use to show comments when clicked
   const [viewLikes, setviewLikes] = useState({}); // use to show likes when clicked
   const [viewCmtBox, setviewCmtBox] = useState({}); // use to show comment box when clicked
   const [page, setPage] = useState(props.page); // use to show comment box when clicked
   const [viewing, setViewing] = useState(props.viewing); // use to show comment box when clicked


   dayjs().format()
   dayjs.extend(relativeTime)

   var timeAtRender = dayjs(Date.now())

  useEffect( () => {
         getCmt().then(cmtz => {
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


   async function getCmt (user2find) { // sends username, password from input, then backend creates s3 bucket in username and stores details on mongo
  
      let data = {Username: viewing}

      let queryType

      if (page === "All") {
         queryType = "GetAllComments"
      } else {
         queryType ="GetUserComments"
      }

      if (user2find) {
         data.Username = sessionUser
      }

      let gqlRequest = "query " + queryType + " ($Username: String!){  " + queryType + " (input: $Username) { Key ID Username Password Email Bio Profpic Photos LastCommentNum Posts { Username SessionUser MainCmt PostNum Time TimeStamp Date Comments { Username Comment Profpic } Likes { Username Profpic } } } }"
      let response = await SendData(gqlRequest, data)
      if ( "errors" in response ){ // if password is a match redirect to profile page
			//{ProcessErrorAlerts("hi", "hi")}
			console.log("Error retrieving user data", response.errors[0].message )
			
		} else {
         return response.data[queryType] 
      }
   }


   async function sendCmt (msgType, cmtAuthr, iD) { // sends comments, replies to comments and likes
      let cmt = "" //in case msgType is reply and comment box not shown so causes null error
      if (msgType === "isCmt") {
         cmt = document.getElementById("cmt").value 
      }   

      if (msgType === "isCmt") {
         let NewCmtInput  = { data: {
            Username : cmtAuthr,
            SessionUser: sessionUser,
            MainCmt : cmt ,
            Time : new Date().toLocaleTimeString('en-GB', { hour: "numeric", minute: "numeric"}), // no longer used
            Date : new Date().toLocaleDateString(), // no longer used
            TimeStamp : Date.now(),    
            ReturnPage: page   ,
            }
         }

         let gqlRequest = "mutation NewComment ($data: SendCmtInput!){ NewComment (input: $data) { Key ID Username Password Email Bio Profpic Photos LastCommentNum Posts { Username SessionUser MainCmt PostNum Time TimeStamp Date Comments { Username Comment Profpic } Likes { Username Profpic } } } }"
         
         SendData(gqlRequest, NewCmtInput).then((response)=> ("errors" in response) ? console.log("error posting data") : setcmt(response.data.NewComment) ); 
      }
      

      if (msgType === "isResponse") {
         let CommentResponse =  { data: {
            AuthorUsername: cmtAuthr ,
            ReplyUsername: sessionUser ,
            ReplyComment:  document.getElementById(iD).value , 
            ReplyProfpic:  "" ,
            PostIndx:   iD ,
            ReturnPage: page
            }
         }

         let gqlRequest = "mutation ReplyComment ($data: ReplyCommentInput!){ ReplyComment (input: $data) { Key ID Username Password Email Bio Profpic Photos LastCommentNum Posts { Username SessionUser MainCmt PostNum Time TimeStamp Date Comments { Username Comment Profpic } Likes { Username Profpic } } } }"
         const reply = await getCmt(sessionUser)
         CommentResponse.data.ReplyProfpic = reply.Profpic
         
         SendData(gqlRequest, CommentResponse).then((response)=> ("errors" in response) ? console.log("error sending response to comment") : setcmt(response.data.ReplyComment))
      }


      if (msgType === "cmtLiked") {
         let SendLikeInput  = { data: {
            Username:   cmtAuthr ,
            LikedBy:   sessionUser ,
            LikeByPic:   "",
            PostIndx:   iD , 
            ReturnPage: page   
            }
         }
         let gqlRequest = "mutation LikeComment ($data: SendLikeInput!){ LikeComment (input: $data) { Key ID Username Password Email Bio Profpic Photos LastCommentNum Posts { Username SessionUser MainCmt PostNum Time TimeStamp Date Comments { Username Comment Profpic } Likes { Username Profpic } } } }"
         getCmt(sessionUser).then((repliersData)=> {SendLikeInput.data.LikeByPic = repliersData.Profpic; SendData(gqlRequest, SendLikeInput).then((response)=> ( "errors" in response) ? console.log("error adding like") : setcmt(response.data.LikeComment) ) ; })
      }

   }

   

   function redirecter () {
      //Navigate("/signIn")
   }


    function goToProfile (){
      //Navigate ("/Profile/" + User + "/home")
    }

    function goToAllPosts (){
      //Navigate ("/profile/" + User + "/all")
    }

    function goToPhotos () {
       if (page !== "media") { // for when your on media tab already
         //Navigate("/profile/" + User + "/media")
       }
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
             <Col style={{}}>
                  <div id="content" className="content content-full-width">
                     <div className="profile">
                        <div className="profile-header">
                           <div className="profile-header-cover"></div>
                              <div className="profile-header-content">
                                 <div className="profile-header-img">
                                    <img calssName="profpics" src={cmt.Profpic} alt=""/>
                                 </div>
                                 <div className="profile-header-info">
                                    <h4 className="m-t-10 m-b-5">{viewing}</h4>
                                    <p className="m-b-10" style={{color:"black"}}>{cmt.Bio? cmt.Bio : "Click the edit profile button to add a bio to your profile now."}</p>
                                    <a href="!" className="btn btn-sm btn-info mb-2" style={{marginRight:"10px"}} onClick={(e)=> {e.preventDefault() ; setPage("editprofile")}}>Edit Profile</a>
                                    <a href="!" className="btn btn-sm btn-info mb-2" style={{marginRight:"10px"}} onClick={(e)=> {e.preventDefault() ; setPage("updatedetails")}}>Update Details</a>
                                    <a href="!" className="btn btn-sm btn-info mb-2" onClick={(e)=> {e.preventDefault(); Logout()}}>Log Out</a>
                                 </div>
                              </div>
                        </div>
                     </div>
                  </div>
                  
             </Col>
             {/* <Col style={{background: "white"}}>2 of 2</Col> */}
           </Row>
           <Row>
             <Col xs lg="3"> 
               <Row>
                   <Col style={{}}>
                     {sessionUser && <button className="login100-form-btn" type="button" onClick={(e) => { e.preventDefault() ; goToProfile()}}>PROFILE</button>}
                   </Col>
               </Row>
               <Row>
                     <Col style={{}}>
                     {sessionUser && <button className="login100-form-btn" type="button"  onClick={(e) => {e.preventDefault(); goToAllPosts()}}>FEED</button>}
                             
                     </Col>
               </Row>
               <Row>
                     <Col style={{}}>
                     {sessionUser && <button className="login100-form-btn" type="button"  onClick={(e) => {e.preventDefault(); goToPhotos()}}>MEDIA</button>}      
                     </Col>
               </Row>
               <Row>
                     <Col>
                     {sessionUser && <button className="login100-form-btn" type="button" onClick={() => {}}>FRIENDS</button>}  
                     </Col>        
               </Row>
             </Col>
             <Col>
                  {sessionUser && page === "media" &&
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
                                       <button className="btn btn-primary f-s-12 rounded-corner" type="button" onClick={() => sendCmt("isCmt", sessionUser, 0)}>Comment</button>
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
                                          <span className="userimage"><img onClick={()=> setViewing(userData.Username)} src={cmt.Profpic} alt=""/></span>
                                       </Col>
                                       <Col>
                                          <Row>
                                             <Col md="auto">
                                                <div className="username" onClick={()=> setViewing(userData.Username)}>{userData.Username}</div>
                                             </Col>
                                             <Col>
                                                <div className="time">{dayjs(userData.TimeStamp).from(timeAtRender) }</div>
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
                                             <span className="fa-stack fa-fw stats-icon">
                                                <i className="fa fa-circle fa-stack-2x text-danger"></i>
                                                <i className="fa fa-heart fa-stack-1x fa-inverse t-plus-1" onClick={()=>sendCmt("cmtLiked", userData.Username, userData.PostNum)}></i>
                                             </span>
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
                                                            <div ClassName = "username" onClick={()=> setViewing(replys.Username)} >{replys.Username}</div> 
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
                                                                  <textarea type="text" className="form-control rounded-corner" id={userData.PostNum} style={{height:"90px", width:"180px"}} placeholder="Reply to the post..."/>
                                                               </div>
                                                            </form>
                                                         </div>
                                                   </Row>
                                                   <Row>
                                                      <span className="input-group-btn p-l-10">
                                                         <button className="btn btn-primary f-s-12 rounded-corner" type="button"  onClick={()=>sendCmt("isResponse", userData.Username, userData.PostNum)}>Comment</button>
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
      {page === "updatedetails" && <UpdateDetails username={sessionUser} /* setSessionUser={setSessionUser()} */ />}
	   {page === "home" &&  <Home sessionuser={sessionUser} page={page} viewing={sessionUser} />}
      {page === "editprofile" && <EditProfile sessionuser={sessionUser} />}
      {page === "signin" && <SignIn/>}
      </>

   )
}

export default memo(Home)