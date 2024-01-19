# I2C

## Dependencies

$ apt-get install gcc-arm-linux-gnueabi

## References
 
 * https://docs.kernel.org/i2c/dev-interface.html
 * i2c-tools, func i2c_smbus_read_i2c_block_data [source](https://github.com/mozilla-b2g/i2c-tools/blob/master/lib/smbus.c)

## Lessons Learned

The NixOS can easily be cross-compiled for ARM. However, apparently it's not possible to copy over and run binaries. This makes NixOS as a dev platform not a good idea. See [link](https://nixos.wiki/wiki/Packaging/Binaries) for more details.