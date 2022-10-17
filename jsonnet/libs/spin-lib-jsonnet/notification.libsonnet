/**
    @file Contains methods for creating Spinnaker Notifications.

    Notifications can be set on stages and pipelines.

    More details on Spinnaker notifications:

    {@link https://spinnaker.io/setup/features/notifications/}
**/

/**
    Creates a Notification object.

    Can be used for stages or pipelines.

    `message` property can be set explicity to provide a custom message for a state.

    @example
        notification.Notification {
            type: notification.Type.slack,
            address: 'my-slack-channel-no-hashtag',
            level: notification.Level.stage,

            when: [
                notification.Level.stage + "." + notification.State.complete
            ],

            message: {
                "stage.complete": {
                    "text": "my stage just finished."
                }
            }
        }

    @property {String} address - Address to which notification will be sent. (eg. slack channel, email address)
    @property {Enum<Level>} level - At which Spinnaker level this notifaction will be sent. See the `Level` enum.
    @property {Enum<Type>} type - Type of notification (eg. slack, email, ...). See the `Type` enum.
    @property {Enum<State>} when - When the notification will be sent (eg. stage failure, stage success, ...). See `State` enum.
**/
local Notification = {
  address: error '`address` (String) property is required for Notification',
  level: error '`level` (String) property is required for Notification',
  type: error '`type` (String) property is required for Notification',
  when: error '`when` (String) property is required for Notification',

  // Other fields:
  // message: {}
};

/**
    A map of stage/pipeline states on which notifications can be sent.

    This map was made from:

    {@link https://github.com/spinnaker/deck/blob/master/app/scripts/modules/core/src/notification/modal/whenOptions.ts}

    @example
        notification.State.complete

    @enum {String}
    @readonly

    @property {String} complete - Complete state of a pipeline/stage.
    @property {String} failed - Failed state of a pipeline/stage.
    @property {String} starting - Starting state of a pipeline/stage.
    @property {String} manualJudgment - Manual Judgement state, when the pipeline reaches that stage. Specific to Manual Judgement stage.
    @property {String} manualJudgmentContinue - When the manual judgement was set to continue. Specific to Manual Judgement stage.
    @property {String} manualJudgmentStop - When the manual judgement was set to stop. Specific to Manual Judgement stage.
**/
local State = {
  complete: 'complete',
  failed: 'failed',
  starting: 'starting',

  //specific to manual judgement stage
  manualJudgment: 'manualJudgment',
  manualJudgmentContinue: 'manualJudgmentContinue',
  manualJudgmentStop: 'manualJudgmentStop',
};

/**
    A map of level for which notifications can be sent.

    This map was made from:

    {@link https://github.com/spinnaker/deck/blob/master/app/scripts/modules/core/src/notification/modal/whenOptions.ts}

    @example
        notification.Level.stage

    @enum {String}
    @readonly

    @property {String} pipeline - Pipeline level type.
    @property {String} stage - Stage level type.
**/
local Level = {
  pipeline: 'pipeline',
  stage: 'stage',
};

/**
    A map of all notification types.

    They must be enabled and configured in the Spinnaker environment to be used.

    This map was made from:

    {@link https://github.com/spinnaker/echo/tree/master/echo-notifications/src/main/groovy/com/netflix/spinnaker/echo/notification}

    Look for `getNotificationType` method in each of the `NotificationAgent`.

    @example
        notification.Type.slack

    @enum {String}
    @readonly

    @property {String} BEARY_CHAT - The notification type for Beary Chat.
    @property {String} DRY_RUN - Dry-run notifications, won't send anything.
    @property {String} EMAIL - The notification type for emailing.
    @property {String} GITHUB_STATUS - The notification type for Github Status.
    @property {String} GOOGLE_CHAT - The notification type for Google Chat.
    @property {String} GOOGLE_CLOUD_BUILD - The notification type for Google Cloud Build.
    @property {String} MICROSOFT_TEAMS - The notification type for Microsoft Teams.
    @property {String} SLACK - The notification type for Slack.
    @property {String} SMS - The notification type for SMS.
**/
local Type = {
  BEARY_CHAT: 'bearychat',
  DRY_RUN: 'dryrun',
  EMAIL: 'email',
  GITHUB_STATUS: 'githubStatus',
  GOOGLE_CHAT: 'googlechat',
  GOOGLE_CLOUD_BUILD: 'googleCloudBuild',
  MICROSOFT_TEAMS: 'microsoftteams',
  SLACK: 'slack',
  SMS: 'sms',
};

/**
    Creates a StateMessage object that ties a message to a state.

    Should be used in conjuction with NewNotification.

    @example
        local myStartingMessage = notification.StateMessage {
            message: "Stage has just started running.",
            state: notification.State.starting,
        }
        local myFailedMessage = notification.StateMessage {
            message: "Stage has Failed.",
            state: notification.State.failed,
        }
        local myCompleteMessage = notification.StateMessage {
            message: "Stage has just completed running.",
            state: notification.State.complete,
        }
        local allOfMyMessages = [
            myStartingMessage,
            myFailedMessage,
            myCompleteMessage
        ]

        ...

        notification.NewNotification(notification.Level.stage, allOfMyMessages) { ... }

    @property {String} message - Message for the notification.
    @property {Enum<State>} state - The state/when this message will be apply (eg. stage failed, stage succeeded). See `State` object.
**/
local StateMessage = {
  message: error '`message` (String) property is required for StateMessage',
  state: error '`state` (String) property is required for StateMessage',
};

/**
    Creates a Notification object with the given messages for a given level.

    Should be used in conjuction with StateMessage.

    @todo Is this a class, or is this a function?

    @example
        local myStartingMessage = notification.StateMessage { ... }
        local myFailedMessage = notification.StateMessage { ... }
        local myCompleteMessage = notification.StateMessage { ... }
        local allOfMyMessages = [
            myStartingMessage,
            myFailedMessage,
            myCompleteMessage
        ]

        ...

        notification.NewNotification {
            type: notification.Type.slack,
            address: 'my-slack-channel-no-hashtag',
            Level:: notification.Level.stage,
            StateMessages:: allOfMyMessages
        }

    @constructs Notification
    @memberof Notification
    @name NewNotification

    @param {Enum<Level>} Level - The level on which to send this notification. See `Level` object.
    @param {Array<StateMessage>} StateMessages - An array of StateMessage objects.

    @return {Notification} A new Notification.
**/
local NewNotification = Notification {
  // Used for formatting.
  local makeWhen(level, stateMessage) = (
    // Usually it's `level.state`, but manualJudgement stages just have it as `state`.
    if std.startsWith(stateMessage.state, State.manualJudgment) then
      '%s' % [stateMessage.state]
    else
      '%s.%s' % [level, stateMessage.state]
  ),

  Level:: error '`Level` (String) property is required for NewNotification',
  StateMessages:: error '`StateMessages` (Array<StateMessage>) property is required for NewNotification',


  level: $.Level,
  message: {
    [makeWhen($.Level, stateMessage)]: {
      text: stateMessage.message,
    }
    for stateMessage in $.StateMessages
  },
  when: [
    makeWhen($.Level, stateMessage)
    for stateMessage in $.StateMessages
  ],
};

/**
    Creates a stage configuration object.

    Contains the notifications and a flag to send the notifications.

    Must be appened to a stage.

    @example
        local myNotification = notification.Notification { ... }
        local stageNotifConfig = StageNotificationsConfiguration {
            notifications = [
                myNotification
            ]
        }

        ...

        local myStage = stage.ManualJudgement {
            ...
        } + stageNotifConfig

    @param {Array<Notification>} notifications - An array of Notification objects.
    @param {Boolean} [sendNotifications=true] - Whether or not to send notifications.
**/
local StageNotificationsConfiguration = {
  notifications: error '`message` (Array<Notification>) notifications is required for StageNotificationsConfiguration',

  sendNotifications: true,
};

// Exposed for public use.
{
  Notification:: Notification,
  Level:: Level,
  State:: State,
  StateMessage:: StateMessage,
  Type:: Type,

  StageNotificationsConfiguration:: StageNotificationsConfiguration,
  NewNotification:: NewNotification,
}
