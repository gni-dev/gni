#pragma once

#include <jni.h>

static inline jclass JNI_FindClass(JNIEnv *env, const char *name)
{
    return (*env)->FindClass(env, name);
}

static inline jthrowable JNI_ExceptionOccurred(JNIEnv *env)
{
    return (*env)->ExceptionOccurred(env);
}

static inline void JNI_ExceptionClear(JNIEnv *env)
{
    (*env)->ExceptionClear(env);
}

static inline jobject JNI_NewGlobalRef(JNIEnv *env, jobject obj)
{
    return (*env)->NewGlobalRef(env, obj);
}

static inline jclass JNI_GetObjectClass(JNIEnv *env, jobject obj)
{
    return (*env)->GetObjectClass(env, obj);
}

static inline jobject JNI_NewObjectA(JNIEnv *env, jclass clazz, jmethodID methodID, const jvalue *args)
{
    return (*env)->NewObjectA(env, clazz, methodID, args);
}

static inline jmethodID JNI_GetMethodID(JNIEnv *env, jclass clazz, const char *name, const char *sig)
{
    return (*env)->GetMethodID(env, clazz, name, sig);
}

static inline jobject JNI_CallObjectMethodA(JNIEnv *env, jobject obj, jmethodID methodID, const jvalue *args)
{
    return (*env)->CallObjectMethodA(env, obj, methodID, args);
}

static inline void JNI_CallVoidMethodA(JNIEnv *env, jobject obj, jmethodID methodID, const jvalue *args)
{
    (*env)->CallVoidMethodA(env, obj, methodID, args);
}

static inline jfieldID JNI_GetStaticFieldID(JNIEnv *env, jclass clazz, const char *name, const char *sig)
{
    return (*env)->GetStaticFieldID(env, clazz, name, sig);
}

static inline jobject JNI_GetStaticObjectField(JNIEnv *env, jclass clazz, jfieldID fieldID)
{
    return (*env)->GetStaticObjectField(env, clazz, fieldID);
}

static inline jstring JNI_NewStringUTF(JNIEnv *env, const char *utf)
{
    return (*env)->NewStringUTF(env, utf);
}

static inline const char *JNI_GetStringUTFChars(JNIEnv *env, jstring string, jboolean *isCopy)
{
    return (*env)->GetStringUTFChars(env, string, isCopy);
}

static inline void JNI_ReleaseStringUTFChars(JNIEnv *env, jstring string, const char *utf)
{
    (*env)->ReleaseStringUTFChars(env, string, utf);
}

static inline jboolean JNI_ExceptionCheck(JNIEnv *env)
{
    return (*env)->ExceptionCheck(env);
}
