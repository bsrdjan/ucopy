# -*- coding: utf-8 -*-

from pyrfc import *
import datetime

# BAPI calls log
def print_log(bapi_return):
    if len(bapi_return) > 0:
        for line in bapi_return:
            print '%s: %s' % (line['TYPE'], line['MESSAGE'])

# Connect to ABAP system
SAPROUTER = '/H/123.12.123.12/E/yt6ntx/H/123.14.131.111/H/'

EC4 = {
    'user'      : 'abapuser',
    'passwd'    : 'abappass',
    'ashost'    : '10.11.12.13',
    'saprouter' : SAPROUTER,
    'sysnr'     : '00',
    'client'    : '300',
    'trace'     : '3',
    'lang'      : 'EN' }

c = Connection(**EC4)

# The sourse user, to be copied
uname_from = 'UNAMEFROM'
# Defaults if source user validity not maintained (undefined)
valid_from = datetime.date(2015,1,19)
valid_to   = datetime.date(2015,12,31)
# New users' password. For automatic generation check CREATE BAPI
initpwd = 'InitPa$$21'

# Users to be created
users= ['UNAMETO1', 'UNAMETO2']

# Get source user details
r = c.call('BAPI_USER_GET_DETAIL', USERNAME = uname_from, CACHE_RESULTS  = ' ')

# Set new users' defaults
if r['LOGONDATA']['GLTGV'] is None:
    r['LOGONDATA']['GLTGV'] = valid_from
if r['LOGONDATA']['GLTGB'] is None:
    r['LOGONDATA']['GLTGB'] = valid_to
password = {'BAPIPWD' : initpwd}

# Create new users
for uname_to in users:
    print uname_to
    r['ADDRESS']['LASTNAME'] = uname_to
    r['ADDRESS']['FULLNAME'] = uname_to

    x = c.call('BAPI_USER_CREATE1',
        USERNAME    = uname_to,
        LOGONDATA   = r['LOGONDATA'],
        PASSWORD    = password,
        DEFAULTS    = r['DEFAULTS'],
        ADDRESS     = r['ADDRESS'],
        COMPANY     = r['COMPANY'],
        REF_USER    = r['REF_USER'],
        PARAMETER   = r['PARAMETER'],
        GROUPS      = r['GROUPS']
    )

    print_log(x['RETURN'])

    x = c.call('BAPI_USER_PROFILES_ASSIGN',
        USERNAME  = uname_to,
        PROFILES  = r['PROFILES']
    )

    print_log(x['RETURN'])

    x = c.call('BAPI_USER_ACTGROUPS_ASSIGN',
        USERNAME       = uname_to,
        ACTIVITYGROUPS = r['ACTIVITYGROUPS']
    )

    print_log(x['RETURN'])

# Finished
print ("%s copied to %d new users\nBye!") % (uname_from, len(users))

