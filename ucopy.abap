REPORT zsbucopy1.

DATA:
  uname_from        LIKE bapibname-bapibname VALUE 'UNAMEFROM',

  ls_logondata      TYPE bapilogond,
  ls_defaults       TYPE bapidefaul,
  ls_address        TYPE bapiaddr3,
  ls_company        TYPE  bapiuscomp,

  lt_parameter      TYPE STANDARD TABLE OF bapiparam,
  lt_profiles       TYPE STANDARD TABLE OF bapiprof,
  lt_activitygroups TYPE STANDARD TABLE OF bapiagr,
  lt_return         TYPE STANDARD TABLE OF bapiret2,
  lt_parameter1     TYPE STANDARD TABLE OF bapiparam1,
  lt_groups         TYPE STANDARD TABLE OF bapigroups,

  uname_to          LIKE bapibname-bapibname VALUE 'UNAMETO',
  is_password       TYPE bapipwd.

is_password-bapipwd = 'Init2014'.

CALL FUNCTION 'BAPI_USER_GET_DETAIL'
  EXPORTING
    username       = uname_from
    cache_results  = ' '
  IMPORTING
    logondata      = ls_logondata
    defaults       = ls_defaults
    address        = ls_address
    company        = ls_company
  TABLES
    parameter      = lt_parameter
    profiles       = lt_profiles
    activitygroups = lt_activitygroups
    return         = lt_return
    parameter1     = lt_parameter1
    groups         = lt_groups.

MOVE uname_to TO: ls_address-lastname, ls_address-fullname.

CALL FUNCTION 'BAPI_USER_CREATE1'
  EXPORTING
    username   = uname_to
    logondata  = ls_logondata
    password   = is_password
    defaults   = ls_defaults
    address    = ls_address
    company    = ls_company
    ref_user   = uname_from
  TABLES
    parameter  = lt_parameter
    return     = lt_return
    groups     = lt_groups
    parameter1 = lt_parameter1.

CALL FUNCTION 'BAPI_USER_PROFILES_ASSIGN'
  EXPORTING
    username = uname_to
  TABLES
    profiles = lt_profiles
    return   = lt_return.

CALL FUNCTION 'BAPI_USER_ACTGROUPS_ASSIGN'
  EXPORTING
    username       = uname_to
  TABLES
    activitygroups = lt_activitygroups
    return         = lt_return.