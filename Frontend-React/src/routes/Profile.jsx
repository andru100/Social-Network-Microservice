import { React, memo } from "react";
import { useState, useEffect } from "react"
import { useNavigate, useParams } from "react-router";
import ChkAuth from './chkAuth';
import SendData from './SendData';

function Home () {
   const [sessionUser, setSessionUser] = useState ("")
   const[cmt, setcmt] = useState({Posts:[], Photos: []}) // holds users data object
   const [viewReply, setviewReply] = useState(false); // use to show comments when clicked
   const [viewLikes, setviewLikes] = useState(false); // use to show likes when clicked
   const Navigate = useNavigate();
   var {User, Page,} = useParams()

   let Today = new Date().toLocaleDateString()

  useEffect( () => {
    ChkAuth().then(user => {
      if (user) {
         setSessionUser(user)
         getCmt().then(cmtz => {
            if (cmtz) {
             setcmt(cmtz)
             console.log("Users data object retrieved is:", cmtz)
            }
         })
      } else {
         setSessionUser(false)
         alert("You need to sign in to view this page")
         Navigate("/signIn")
      }
    })
  },[Navigate]);


   async function getCmt (user2find) { // sends username, password from input, then backend creates s3 bucket in username and stores details on mongo
  
      let data = {Username: User}

      let queryType

      if (Page === "All") {
         queryType = "GetAllComments"
      } else {
         queryType ="GetUserComments"
      }

      if (user2find) {
         data.Username = sessionUser
      }

      let gqlRequest = "query " + queryType + " ($Username: String!){  " + queryType + " (input: $Username) { Key ID Username Password Email Bio Profpic Photos LastCommentNum Posts { Username SessionUser MainCmt PostNum Time TimeStamp Date Comments { Username Comment Profpic } Likes { Username Profpic } } } }"
      let response = await SendData(gqlRequest, data)
      if (response) {
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
            Time : new Date().toLocaleTimeString('en-GB', { hour: "numeric", minute: "numeric"}),
            Date : Today,
            TimeStamp : Date.now(),    
            ReturnPage: Page   ,
            }
         }

         let gqlRequest = "mutation NewComment ($data: SendCmtInput!){ NewComment (input: $data) { Key ID Username Password Email Bio Profpic Photos LastCommentNum Posts { Username SessionUser MainCmt PostNum Time TimeStamp Date Comments { Username Comment Profpic } Likes { Username Profpic } } } }"
         
         SendData(gqlRequest, NewCmtInput).then((response)=> response ? setcmt(response.data.NewComment) :  console.log("error posting new comment") ); 
      }
      

      if (msgType === "isResponse") {
         let CommentResponse =  { data: {
            AuthorUsername: cmtAuthr ,
            ReplyUsername: sessionUser ,
            ReplyComment:  document.getElementById(iD).value , 
            ReplyProfpic:  "" ,
            PostIndx:   iD ,
            ReturnPage: Page
            }
         }

         let gqlRequest = "mutation ReplyComment ($data: ReplyCommentInput!){ ReplyComment (input: $data) { Key ID Username Password Email Bio Profpic Photos LastCommentNum Posts { Username SessionUser MainCmt PostNum Time TimeStamp Date Comments { Username Comment Profpic } Likes { Username Profpic } } } }"
         const reply = await getCmt(sessionUser)
         CommentResponse.data.ReplyProfpic = reply.Profpic
         
         SendData(gqlRequest, CommentResponse).then((response)=> response ? setcmt(response.data.ReplyComment) :  console.log("error sending response to comment") )
      }


      if (msgType === "cmtLiked") {
         let SendLikeInput  = { data: {
            Username:   cmtAuthr ,
            LikedBy:   sessionUser ,
            LikeByPic:   "",
            PostIndx:   iD , 
            ReturnPage: Page   
            }
         }
         let gqlRequest = "mutation LikeComment ($data: SendLikeInput!){ LikeComment (input: $data) { Key ID Username Password Email Bio Profpic Photos LastCommentNum Posts { Username SessionUser MainCmt PostNum Time TimeStamp Date Comments { Username Comment Profpic } Likes { Username Profpic } } } }"
         getCmt(sessionUser).then((repliersData)=> {SendLikeInput.data.LikeByPic = repliersData.Profpic; SendData(gqlRequest, SendLikeInput).then((response)=> response ? setcmt(response.data.LikeComment) :  console.log("error when sending comment like") ); })
      }

   }

   function redirecter () {
      Navigate("/signIn")
   }

   function EditProfile (){
      Navigate ("/editProfile/"+sessionUser)
    }

    function goToProfile (){
      Navigate ("/Profile/" + User + "/home")
    }

    function goToAllPosts (){
      Navigate ("/profile/" + User + "/all")
    }

    function goToPhotos () {
       if (Page !== "media") { // for when your on media tab already
         Navigate("/profile/" + User + "/media")
       }
    }

    function Logout (){
      localStorage.removeItem('jwt_token')
      Navigate ("/")
    }



return (
       
<div className="container">
   <div className="row">
      <div className="col-md-12">
         <div id="content" className="content content-full-width">
            {/* begin profile */}
            <div className="profile">
               <div className="profile-header">
                  <div className="profile-header-cover"></div>
                  <div className="profile-header-content">
                     <div className="profile-header-img">
                        <img calssName="profpics" src={cmt.Profpic} alt=""/>
                     </div>
                     <div className="profile-header-info">
                        <h4 className="m-t-10 m-b-5">{User}</h4>
                        <p className="m-b-10" style={{color:"black"}}>{cmt.Bio? cmt.Bio : "Click the edit profile button to add a bio to your profile now."}</p>
                        <a href="!" className="btn btn-sm btn-info mb-2" style={{marginRight:"10px"}} onClick={(e)=> {e.preventDefault() ; EditProfile()}}>Edit Profile</a>
                        <a href="!" className="btn btn-sm btn-info mb-2" onClick={(e)=> {e.preventDefault(); Logout()}}>Log Out</a>
                     </div>
                  </div>
                  {/* Navigation tabs */}
                  <ul className="profile-header-tab nav nav-tabs">
                      {Page === "x" ? <li className="nav-item"><a href="!" className="nav-link active show" data-toggle="tab" onClick={(e) => { e.preventDefault() ; goToProfile()}}>PROFILE</a></li>: <li className="nav-item"><a href="!" className="nav-link" data-toggle="tab" onClick={(e) => {e.preventDefault(); goToProfile()}}>PROFILE</a></li>} 
                      {Page === "all" ? <li className="nav-item"><a href="!" className="nav-link active show" data-toggle="tab" onClick={(e) => {e.preventDefault(); goToAllPosts()}}>FEED</a></li>: <li className="nav-item"><a href="!" className="nav-link" data-toggle="tab" onClick={(e) => {e.preventDefault() ; goToAllPosts()}}>FEED</a></li>} 
                      {Page === "media" ? <li className="nav-item"><a href="!" className="nav-link active show" data-toggle="tab" onClick={(e) => {e.preventDefault(); goToPhotos()}}>MEDIA</a></li>: <li className="nav-item"><a href="!" className="nav-link" data-toggle="tab" onClick={(e) => {e.preventDefault() ; goToPhotos()}}>MEDIA</a></li>} 
                      {Page === "friends" ? <li className="nav-item"><a href="!" className="nav-link active show" data-toggle="tab">MEDIA</a></li>: <li className="nav-item"><a href="!" className="nav-link" data-toggle="tab" onClick={() => {}}>FRIENDS</a></li>} 
                  </ul>
               </div>
            </div>
            {/*show media page */}
            {sessionUser && Page === "media" ? 
            <div className="connected-container">
               <div className="gif-grid">
                  {cmt.Photos.map((pic) => (
                     <div className="gif-item" key={pic}>
                        <img src={pic} alt={pic} />
                     </div>
                  ))}
               </div>
            </div>
            : 
            <div></div>}
            {/*render users profile page or news feed showing all comments*/}
            {(sessionUser === User && Page !== "media") || (sessionUser && Page === "all") ? <div className="input">
                     <form action="">
                        <div className="input-group">
                           <textarea style={{height: "36px"}} type="text" className="form-control rounded-corner" id="cmt" placeholder="Write a comment..."/>
                           <span className="input-group-btn p-l-10">
                              <button className="btn btn-primary f-s-12 rounded-corner" type="button" onClick={() => sendCmt("isCmt", sessionUser, 0)}>Comment</button>
                           </span>
                        </div>
                  </form>
              </div> : <div></div>}
            <div className="profile-content">
               <div className="tab-content p-0">
                  <div className="tab-pane fade active show" id="profile-post">
                     <ul className="timeline">
                     { Page !== "media" ? 
                        cmt.Posts.map((userData)=> (
                           <li>
                           <div className="timeline-time">
                              <span className="date">{userData.Date === Today ? "Today" : userData.Date}</span>
                              <span className="time">{userData.Time}</span>
                           </div>
                           <div className="timeline-icon">
                              <a href="!">&nbsp;</a>
                           </div>
                           <div className="timeline-body">
                              <div className="timeline-header">
                                 <span className="userimage"><img src={cmt.Profpic} alt=""/></span>
                                 <a className="username" href= {process.env.REACT_APP_FRONTEND + "/Profile/" + userData.Username + "/home"}>{userData.Username}</a>
                              </div>
                              <div className="timeline-content">
                                 <p>
                                    {userData.MainCmt}
                                 </p>
                              </div>
                              <div className="timeline-likes">
                                 <div className="stats-left">
                                    <span className="fa-stack fa-fw stats-icon">
                                    <i className="fa fa-circle fa-stack-2x text-danger"></i>
                                    <i className="fa fa-heart fa-stack-1x fa-inverse t-plus-1" onClick={()=>sendCmt("cmtLiked", userData.Username, userData.PostNum)}></i>
                                    </span>
                                    <span className="stats-text" onClick={() => {viewReply && setviewReply(!viewReply) ; setviewLikes(!viewLikes)}}>{userData.Likes?.length} Likes</span>
                                    <span className="stats-text" onClick={() => {viewLikes && setviewLikes(!viewLikes) ; setviewReply(!viewReply)}}>{userData.Comments?.length} Comments</span>
                                    {viewLikes &&
                                      userData.Likes.map((Likes)=> (
                                          <div className="timeline-header">
                                          {<span className="userimage"><img src={Likes.Profpic} alt=""/></span>  }
                                          <a ClassName = "username" href= {process.env.REACT_APP_FRONTEND + "/Profile/" + Likes.Username + "/home"}>{Likes.Username}</a> 
                                          </div> 
                                      ))
                                    } 
                                    {viewReply &&
                                      userData.Comments.map((replys)=> (
                                       <div className="timeline-header">
                                          {<span className="userimage"><img src={replys.Profpic} alt=""/></span> }
                                          <span ><a href= {process.env.REACT_APP_FRONTEND + "/Profile/" + replys.Username + "/home"}>{replys.Username}</a> <small></small></span>
                                          <div>
                                             <span>{replys.Comment}</span>
                                          </div>
                                       </div>
                                      ))
                                    } 
                                 </div>
                                 <div className="stats">
                                 </div>
                              </div>
                              <div className="timeline-comment-box">
                                 <div className="input">
                                    <form action="">
                                       {sessionUser ? 
                                          <div className="input-group">
                                          <textarea type="text" className="form-control rounded-corner" id={userData.PostNum} style={{height:"90px", width:"180px"}} placeholder="Reply to the post..."/>
                                          <div style={{height:"10px", width:"180px"}}></div>
                                          <span className="input-group-btn p-l-10">
                                             <button className="btn btn-primary f-s-12 rounded-corner" type="button"  onClick={()=>sendCmt("isResponse", userData.Username, userData.PostNum)}>Comment</button>
                                          </span>
                                          </div>
                                          :
                                          <div className="input-group">
                                          <input type="text" className="form-control rounded-corner" id={userData.PostNum} placeholder="Write a comment..."/>
                                          <span className="input-group-btn p-l-10">
                                             <button className="btn btn-primary f-s-12 rounded-corner" type="button" onClick={()=>redirecter()} >Sign In to Commment</button>
                                          </span>
                                          </div>}
                                    </form>
                                 </div>
                              </div>
                           </div>
                        </li>
                        ))
                        : 
                         <div></div>
                     }
                     </ul>
                  </div>
               </div>
            </div>
         </div>
      </div>
   </div>
</div>
    )
}

export default memo(Home)