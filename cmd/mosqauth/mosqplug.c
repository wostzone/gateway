// Credit: http://github.com/iegomez/mosquitto-go-auth
// This is a copy of the C file to receive mosquitto calls and pass them on to go
// It requires mosquitto-dev and libmosquitto-dev packages to be installed
#include <string.h>
#include <stdlib.h>
#include <stdio.h>
#include <errno.h>

#include <mosquitto_broker.h>
#include <mosquitto_plugin.h>
#include <mosquitto.h>
#include <openssl/x509.h>

#define mosquitto_auth_opt mosquitto_opt

// main.h is generated when main.go is compiled
#include "main.h"

// Same constant as one in go-auth.go.
// #define AuthRejected 0
// #define AuthGranted 1
// #define AuthError 2

int mosquitto_auth_plugin_version(void)
{
#ifdef MOSQ_AUTH_PLUGIN_VERSION
#if MOSQ_AUTH_PLUGIN_VERSION == 5
  return 4; // This is v2.0, use the backwards compatibility
#else
  return MOSQ_AUTH_PLUGIN_VERSION;
#endif
#else
  return 4;
#endif
}

int mosquitto_auth_plugin_init(void **user_data, struct mosquitto_auth_opt *auth_opts, int auth_opt_count)
{
  /*
    Pass auth_opts hash as keys and values char* arrays to Go in order to initialize them there.
  */
  GoInt32 opts_count = auth_opt_count;

  GoString keys[auth_opt_count];
  GoString values[auth_opt_count];
  int i;
  struct mosquitto_auth_opt *o;
  for (i = 0, o = auth_opts; i < auth_opt_count; i++, o++)
  {
    GoString opt_key = {o->key, strlen(o->key)};
    GoString opt_value = {o->value, strlen(o->value)};
    keys[i] = opt_key;
    values[i] = opt_value;
  }

  GoSlice keysSlice = {keys, auth_opt_count, auth_opt_count};
  GoSlice valuesSlice = {values, auth_opt_count, auth_opt_count};

  AuthPluginInit(keysSlice, valuesSlice, opts_count);
  return MOSQ_ERR_SUCCESS;
}

int mosquitto_auth_plugin_cleanup(void *user_data, struct mosquitto_auth_opt *auth_opts, int auth_opt_count)
{
  AuthPluginCleanup();
  return MOSQ_ERR_SUCCESS;
}

int mosquitto_auth_security_init(void *user_data, struct mosquitto_auth_opt *auth_opts, int auth_opt_count, bool reload)
{
  return MOSQ_ERR_SUCCESS;
}

int mosquitto_auth_security_cleanup(void *user_data, struct mosquitto_auth_opt *auth_opts, int auth_opt_count, bool reload)
{
  return MOSQ_ERR_SUCCESS;
}

#if MOSQ_AUTH_PLUGIN_VERSION >= 4
int mosquitto_auth_unpwd_check(void *user_data, struct mosquitto *client, const char *username, const char *password)
#else
int mosquitto_auth_unpwd_check(void *userdata, const struct mosquitto *client, const char *username, const char *password)
#endif
{
  const char *clientid = mosquitto_client_id(client);
  const char *ip = mosquitto_client_address(client);

  if (username == NULL || password == NULL)
  {
    printf("error: received null username or password for unpwd check\n");
    fflush(stdout);
    return MOSQ_ERR_AUTH;
  }

  GoString go_username = {username, strlen(username)};
  GoString go_password = {password, strlen(password)};
  GoString go_clientid = {clientid, strlen(clientid)};
  GoString go_clientip = {clientid, strlen(ip)};

  GoUint8 ret = AuthUnpwdCheck(go_clientid, go_username, go_password, go_clientip);
  return ret;
}

#if MOSQ_AUTH_PLUGIN_VERSION >= 4
int mosquitto_auth_acl_check(void *user_data, int access, struct mosquitto *client, const struct mosquitto_acl_msg *msg)
#else
int mosquitto_auth_acl_check(void *userdata, int access, const struct mosquitto *client, const struct mosquitto_acl_msg *msg)
#endif
{
  const char *clientid = mosquitto_client_id(client);
  const char *username = mosquitto_client_username(client);
  const char *topic = msg->topic;
  const char *ip = mosquitto_client_address(client);

  X509 *cert = mosquitto_client_certificate(client); // client uses cert auth
  const char *subjname = "";
  if (cert != NULL)
  {
    subjname = X509_NAME_oneline(X509_get_subject_name(cert), NULL, 0);
  }

  if (clientid == NULL || username == NULL || topic == NULL || access < 1)
  {
    printf("error: received null username, clientid or topic, or access is equal or less than 0 for acl check\n");
    fflush(stdout);
    return MOSQ_ERR_ACL_DENIED;
  }

  GoInt32 go_access = access;
  GoString go_clientid = {clientid, strlen(clientid)};
  GoString go_username = {username, strlen(username)};
  GoString go_topic = {topic, strlen(topic)};
  GoString go_subjname = {subjname, strlen(subjname)};

  // ssl testing
  printf("mosquitto_auth_acl_check: clientId=%s, subjectname=%s\n", clientid, subjname);
  fflush(stdout);
  //---
  GoUint8 ret = AuthAclCheck(go_clientid, go_username, go_topic, go_access, go_subjname);
  X509_free(cert);
  return ret;
}

#if MOSQ_AUTH_PLUGIN_VERSION >= 4
int mosquitto_auth_psk_key_get(void *user_data, struct mosquitto *client, const char *hint, const char *identity, char *key, int max_key_len)
#else
int mosquitto_auth_psk_key_get(void *userdata, const struct mosquitto *client, const char *hint, const char *identity, char *key, int max_key_len)
#endif
{
  return MOSQ_ERR_AUTH;
}