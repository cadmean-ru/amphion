// Objective-C API for talking to github.com/cadmean-ru/amphion/common/dispatch Go package.
//   gobind -lang=objc github.com/cadmean-ru/amphion/common/dispatch
//
// File is generated by gobind. Do not edit.

#ifndef __Dispatch_H__
#define __Dispatch_H__

@import Foundation;
#include "ref.h"
#include "Universe.objc.h"


@class DispatchLooperImpl;
@class DispatchMessage;
@class DispatchMessageQueue;
@class DispatchWorkItemFunc;
@protocol DispatchLooper;
@class DispatchLooper;
@protocol DispatchMessageDispatcher;
@class DispatchMessageDispatcher;
@protocol DispatchMessageHandler;
@class DispatchMessageHandler;
@protocol DispatchWorkDispatcher;
@class DispatchWorkDispatcher;
@protocol DispatchWorkItem;
@class DispatchWorkItem;

@protocol DispatchLooper <NSObject>
- (id<DispatchMessageDispatcher> _Nullable)getMessageDispatcher;
- (void)loop;
@end

@protocol DispatchMessageDispatcher <NSObject>
- (void)sendMessage:(DispatchMessage* _Nullable)message;
@end

@protocol DispatchMessageHandler <NSObject>
- (void)onMessage:(DispatchMessage* _Nullable)msg;
@end

@protocol DispatchWorkDispatcher <NSObject>
- (void)execute:(id<DispatchWorkItem> _Nullable)item;
@end

@protocol DispatchWorkItem <NSObject>
- (void)execute;
@end

@interface DispatchLooperImpl : NSObject <goSeqRefInterface, DispatchLooper, DispatchMessageDispatcher, DispatchWorkDispatcher> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
// skipped constructor LooperImpl.NewLooperImpl with unsupported parameter or return types

- (nullable instancetype)initCompat:(long)messageQueueBuffer;
- (void)execute:(id<DispatchWorkItem> _Nullable)item;
- (id<DispatchMessageDispatcher> _Nullable)getMessageDispatcher;
- (id<DispatchWorkDispatcher> _Nullable)getWorkDispatcher;
- (void)loop;
- (void)removeMessageHandler:(long)what;
- (void)sendMessage:(DispatchMessage* _Nullable)message;
// skipped method LooperImpl.SetMessageHandler with unsupported parameter or return types

@end

/**
 * Message holds an int message code and some data.
 */
@interface DispatchMessage : NSObject <goSeqRefInterface> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (nullable instancetype)init:(long)what;
// skipped constructor Message.NewMessageFrom with unsupported parameter or return types

// skipped constructor Message.NewMessageFromWithAnyData with unsupported parameter or return types

// skipped constructor Message.NewMessageWithAnyData with unsupported parameter or return types

- (nullable instancetype)initWithStringData:(long)what data:(NSString* _Nullable)data;
@property (nonatomic) long id_;
@property (nonatomic) long what;
@property (nonatomic) NSString* _Nonnull strData;
// skipped field Message.AnyData with unsupported type: interface{}

// skipped field Message.Sender with unsupported type: interface{}

- (NSString* _Nonnull)string;
@end

/**
 * MessageQueue is a thread-safe message buffer.
In general, it can be used to communicate between different threads.
It is mostly used to communicate between the engine and the frontend.
There are two channels under the hood - main and secondary.
The size of main channel is specified in the constructor function,
whereas the size of the secondary channel equals half the size of the main channel.
You can lock the main channel, so new messages won't be sent to it. Then all messages in it can be processed.
While the main channel is locked, messages are sent to the secondary channel.
As soon as the main channel is unlocked, all messages from the secondary channel are sent to main channel.
 */
@interface DispatchMessageQueue : NSObject <goSeqRefInterface> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
// skipped constructor MessageQueue.NewMessageQueue with unsupported parameter or return types

/**
 * Close closes the underlying channels.
 */
- (void)close;
/**
 * Dequeue Removes the first message from the queue.
If queue is empty returns an empty Message.
 */
- (DispatchMessage* _Nullable)dequeue;
/**
 * DequeueBlocking Removes the first message from the queue.
If the queue is empty waits for a message to come.
 */
- (DispatchMessage* _Nullable)dequeueBlocking;
/**
 * Enqueue sends the specified message to the end of the queue.
Depending on the current state of the MessageQueue it sent to either main or secondary channel.
 */
- (void)enqueue:(DispatchMessage* _Nullable)message;
// skipped method MessageQueue.GetSize with unsupported parameter or return types

/**
 * IsEmpty indicate if the queue is empty, i.e. size == 0.
 */
- (BOOL)isEmpty;
/**
 * LockMainChannel locks the main channel.
While it is locked new messages will be sent to the secondary channel.
The locked state remains until UnlockMainChannel is called.
 */
- (void)lockMainChannel;
/**
 * UnlockMainChannel unlocks the main channel.
All messages from the secondary channel are sent to the main one.
After that new messages are sent to the main channel again.
 */
- (void)unlockMainChannel;
/**
 * WaitForMessage blocks the calling thread until a message with the given code is read from the queue.
 */
- (void)waitForMessage:(long)what;
@end

@interface DispatchWorkItemFunc : NSObject <goSeqRefInterface, DispatchWorkItem> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
// skipped constructor WorkItemFunc.NewWorkItemFunc with unsupported parameter or return types

- (void)execute;
@end

FOUNDATION_EXPORT const int64_t DispatchMessageEmpty;
FOUNDATION_EXPORT const int64_t DispatchMessageExit;
FOUNDATION_EXPORT const int64_t DispatchMessageWorkExec;

// skipped function NewLooperImpl with unsupported parameter or return types


FOUNDATION_EXPORT DispatchLooperImpl* _Nullable DispatchNewLooperImplCompat(long messageQueueBuffer);

FOUNDATION_EXPORT DispatchMessage* _Nullable DispatchNewMessage(long what);

// skipped function NewMessageFrom with unsupported parameter or return types


// skipped function NewMessageFromWithAnyData with unsupported parameter or return types


// skipped function NewMessageQueue with unsupported parameter or return types


// skipped function NewMessageWithAnyData with unsupported parameter or return types


FOUNDATION_EXPORT DispatchMessage* _Nullable DispatchNewMessageWithStringData(long what, NSString* _Nullable data);

// skipped function NewWorkItemFunc with unsupported parameter or return types


@class DispatchLooper;

@class DispatchMessageDispatcher;

@class DispatchMessageHandler;

@class DispatchWorkDispatcher;

@class DispatchWorkItem;

@interface DispatchLooper : NSObject <goSeqRefInterface, DispatchLooper> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (id<DispatchMessageDispatcher> _Nullable)getMessageDispatcher;
- (void)loop;
@end

@interface DispatchMessageDispatcher : NSObject <goSeqRefInterface, DispatchMessageDispatcher> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (void)sendMessage:(DispatchMessage* _Nullable)message;
@end

@interface DispatchMessageHandler : NSObject <goSeqRefInterface, DispatchMessageHandler> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (void)onMessage:(DispatchMessage* _Nullable)msg;
@end

@interface DispatchWorkDispatcher : NSObject <goSeqRefInterface, DispatchWorkDispatcher> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (void)execute:(id<DispatchWorkItem> _Nullable)item;
@end

@interface DispatchWorkItem : NSObject <goSeqRefInterface, DispatchWorkItem> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (void)execute;
@end

#endif
