# TODO:

Error handling []
Set up mysql and configure it for multiple agent handlings []
Set up tasking module / Database []
Fix some of the verbiage on the modules so that it makes a bit more sense []
Maybe some autocomplete and up arrow stuff []
Add the ability to upload to other albums []

Change all the options so that when you type "options" and the value exists, it queries the database and not the global maps []
		--> Eh, I think this is fine, I'll try and see what others think



# DONE:
Display walkthroughs for each module [√]
Fix the tasking module since I commented out the upload function [√]
Error handling for existing tables and what not [√]








Thoughts right now is to have a target upload an encoded image after getting tasking, just not sure if the agent definition should be on response or on tasking

How this should work:
Set up a tasking image for a client -> client grabs it, runs it, uploads new(?) image based on description -> server (C2) is keeping track of each agent and checking for responses (New album for each agent/tasking)


 My thought right now is that this will look for any descriptions that mention the word "response"
 For right now, and then maybe get a big dict of words that will mean specific things

 Beta build:

 1. Create an image with an encoded command in it
 2. Create an album for this to go in, with a particular title (That will be created within the payload so that when the agent executes the payload, it will know what album to look for)
 3. Add "tasking image" to this album (or any other /shrug)
 4. <Target runs payload, grabs image and decodes/runs payload, uploads new image with response to the same album>
 5. C2 Server will pull any album it has marked as "ACTIVE", and see if there are any new images
 6. Either alert the operator or have to operator do a manual check, and show that there is a new image with a reponse in it
 7. Either auto-show the response to the operator or have the operator use the Reponse module to check the response for a particular target



