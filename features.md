
# Features

X = Done

P = In progress

C = Current task

X 1. From the command line, I can log a [280 character] message. Save to s3 bucket or local file.

    * This could be improved to not require quotes around characters that need
      to be escaped (like ' or ?). Or just documentation of this fact.

X 2. I can encrypt the message with a key I set.

3. I can change my encryption key if I know the old one.

X 4. I can read all my messages.

    * Should handle lots of messages, and list limit from AWS/Cloud

5. I can read the last N messages.

    * Make N (10?) the default?

X 6. I can view the day/time I wrote a message.

X 7. I can delete a message.

X 8. I can search messages by content.

    * Search is "dumb" right now, and just does a text match. No indexing. Could
      be better!

9. I can search messages by date.

P 10. All commands are tested

   * Storage needs tests
    
11. Cobra defaults are set properly
    
    * Specifically the config file - this should be the file we're creating
      already for storage, if any.

12. I can use multiple cloud environments
    
X 13. Setup Actions to run CI and build/publish dockerfile
