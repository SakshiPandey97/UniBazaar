const LOGGING_ENABLED=true

export const logger=(message,args)=>{
    if(LOGGING_ENABLED){
        console.log(message+" "+args)
    }
}