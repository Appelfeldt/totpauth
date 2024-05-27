# totpauth
totpauth is a CLI tool that generates time-based one-time passwords using a keycode.

## Synopsis
```totpauth <filepath> [flags]```

Example:  
```totpauth ./secret.key```  
```cat ./secret.key | totpauth```
      
## Flags
```--timestart    (Default value: 0)```  
   Sets the Unix time from which to start counting time steps.
   
```--timestep    (Default value: 30)```  
   Determine the interval at which new otp passwords are generated.
