#include <linux/init.h>
#include <linux/module.h>

MODULE_AUTHOR("Paul R. Tagliamonte");
MODULE_DESCRIPTION("Test driver");
MODULE_LICENSE("MIT");

static int __init test_init(void) { return 0; }
static void __exit test_exit(void) {}

module_init(test_init);
module_exit(test_exit);
