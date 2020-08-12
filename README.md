# nest-logs

Useful for forwarding Cloudwatch Logs to a Kinesis Stream.

If you're logGroup name is of the form `stack/app/STAGE` the appropriate tags
will be merged into each message. Note, you will need to be logging JSON
messages for this to work.

To deploy, first build the zip:

    $ ./build.sh

(You can also grab a latest release from Github and Zip that yourself if you
don't want to install Go.)

Then manually create the lambda. You can easily add the Trigger(s) for your
desired log groups in the Console. You'll also want to add permission to
`PutRecords` on the relevant Kinesis Stream.

Nb: you'll need to make sure the `handler` for the Lambda is set to `main`, and
select `Go 1.x` as the runtime. 256mb memory is recommended.
