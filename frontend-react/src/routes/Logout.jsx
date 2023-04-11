export default function Logout (){

   function DeleteSession (){
      localStorage.removeItem('jwt_token')
   }

   return (
      <a href="!" className="btn btn-sm btn-info mb-2" onClick={(e)=> {e.preventDefault(); DeleteSession()}}>Log Out</a>
                           
   )
 }