# include "scop.h"


const uint32_t WIDTH = 800;
const uint32_t HEIGHT = 600;

const std::vector<const char*> validationLayers = {
    "VK_LAYER_KHRONOS_validation"
};

const std::vector<const char*> deviceExtensions = {
    VK_KHR_SWAPCHAIN_EXTENSION_NAME
};


// #ifdef NDEBUG
    // const bool enableValidationLayers = false;
// #else
const bool enableValidationLayers = true;
// #endif


std::vector<const char*> getRequiredExtensions() {
    uint32_t glfwExtensionCount = 0;
    const char** glfwExtensions;
    glfwExtensions = glfwGetRequiredInstanceExtensions(&glfwExtensionCount);

    std::vector<const char*> extensions(glfwExtensions, glfwExtensions + glfwExtensionCount);

    if (enableValidationLayers) {
        extensions.push_back(VK_EXT_DEBUG_UTILS_EXTENSION_NAME);
    }

    return extensions;
}

bool checkValidationLayerSupport() {
    uint32_t layerCount;
    vkEnumerateInstanceLayerProperties(&layerCount, nullptr);

    std::vector<VkLayerProperties> availableLayers(layerCount);
    vkEnumerateInstanceLayerProperties(&layerCount, availableLayers.data());

    for (const char* layerName : validationLayers) {
        bool layerFound = false;
        for (const auto& layerProperties : availableLayers) {
            if (strcmp(layerName, layerProperties.layerName) == 0) {
                layerFound = true;
                break;
            }
        }
        if (!layerFound) 
            return false;
    }
    return true;
}

void enumerateVkExtension() {
    uint32_t extensionCount = 0;
    vkEnumerateInstanceExtensionProperties(nullptr, &extensionCount, nullptr);

    std::vector<VkExtensionProperties> extensions(extensionCount);

    vkEnumerateInstanceExtensionProperties(nullptr, &extensionCount, extensions.data());
    std::cout << "available extensions:\n";

    for (const auto& extension : extensions) {
        std::cout << '\t' << extension.extensionName << '\n';
    }
}

static VKAPI_ATTR VkBool32 VKAPI_CALL debugCallback(
    VkDebugUtilsMessageSeverityFlagBitsEXT messageSeverity,
    VkDebugUtilsMessageTypeFlagsEXT messageType,
    const VkDebugUtilsMessengerCallbackDataEXT* pCallbackData,
    void* pUserData) {

    (void)messageType;
    (void)pUserData;
    (void)messageSeverity;

    if (messageSeverity >= VK_DEBUG_UTILS_MESSAGE_SEVERITY_WARNING_BIT_EXT)
        std::cerr << "validation layer: " << pCallbackData->pMessage << std::endl;

    return VK_FALSE;
}

VkResult CreateDebugUtilsMessengerEXT(VkInstance instance, const VkDebugUtilsMessengerCreateInfoEXT* pCreateInfo, const VkAllocationCallbacks* pAllocator, VkDebugUtilsMessengerEXT* pDebugMessenger) {
    auto func = (PFN_vkCreateDebugUtilsMessengerEXT) vkGetInstanceProcAddr(instance, "vkCreateDebugUtilsMessengerEXT");
    if (func != nullptr) {
        return func(instance, pCreateInfo, pAllocator, pDebugMessenger);
    } else {
        return VK_ERROR_EXTENSION_NOT_PRESENT;
    }
}

void DestroyDebugUtilsMessengerEXT(VkInstance instance, VkDebugUtilsMessengerEXT debugMessenger, const VkAllocationCallbacks* pAllocator) {
    auto func = (PFN_vkDestroyDebugUtilsMessengerEXT) vkGetInstanceProcAddr(instance, "vkDestroyDebugUtilsMessengerEXT");
    if (func != nullptr) {
        func(instance, debugMessenger, pAllocator);
    }
}

void populateDebugMessengerCreateInfo(VkDebugUtilsMessengerCreateInfoEXT& createInfo) {
    createInfo = {};
    createInfo.sType = VK_STRUCTURE_TYPE_DEBUG_UTILS_MESSENGER_CREATE_INFO_EXT;
    createInfo.messageSeverity = VK_DEBUG_UTILS_MESSAGE_SEVERITY_VERBOSE_BIT_EXT | VK_DEBUG_UTILS_MESSAGE_SEVERITY_WARNING_BIT_EXT | VK_DEBUG_UTILS_MESSAGE_SEVERITY_ERROR_BIT_EXT;
    createInfo.messageType = VK_DEBUG_UTILS_MESSAGE_TYPE_GENERAL_BIT_EXT | VK_DEBUG_UTILS_MESSAGE_TYPE_VALIDATION_BIT_EXT | VK_DEBUG_UTILS_MESSAGE_TYPE_PERFORMANCE_BIT_EXT;
    createInfo.pfnUserCallback = debugCallback;
}



struct QueueFamilyIndices {
    std::optional<uint32_t> graphicsFamily;
    std::optional<uint32_t> presentFamily;
    
    bool isComplete() {
        return graphicsFamily.has_value() && presentFamily.has_value();
    }
};

struct SwapChainSupportDetails {
    VkSurfaceCapabilitiesKHR capabilities;
    std::vector<VkSurfaceFormatKHR> formats;
    std::vector<VkPresentModeKHR> presentModes;
};


static std::vector<char> readFile(const std::string& filename) {
    std::ifstream file(filename, std::ios::ate | std::ios::binary);

    if (!file.is_open()) {
        throw std::runtime_error("failed to open file!");
    }

    size_t fileSize = (size_t) file.tellg();
    std::vector<char> buffer(fileSize);
    file.seekg(0);
    file.read(buffer.data(), fileSize);
    file.close();
    std::cout << "shader (" << filename << ") fileSize: " << buffer.size() << std::endl;
    return buffer;
}

class TApp 
{
    public:
        void run()
        {
            initWindow();
            initVulkan();
            mainLoop();
            cleanup();
        }

    private:
        GLFWwindow*                 window;
        VkInstance                  instance;
        VkDebugUtilsMessengerEXT    debugMessenger;
        VkSurfaceKHR                surface;
        VkPhysicalDevice            physicalDevice = VK_NULL_HANDLE;
        VkDevice                    device;
        VkQueue                     graphicsQueue;
        VkQueue                     presentQueue;

        VkSwapchainKHR              swapChain;
        VkFormat                    swapChainImageFormat;
        VkExtent2D                  swapChainExtent;
        std::vector<VkImage>        swapChainImages;
        std::vector<VkImageView>    swapChainImageViews;

        VkRenderPass                renderPass;
        VkPipelineLayout            pipelineLayout;
        VkPipeline                  graphicsPipeline;

        std::vector<VkFramebuffer>  swapChainFramebuffers;

        VkCommandPool               commandPool;
        VkCommandBuffer             commandBuffer;

// utils
        QueueFamilyIndices findQueueFamilies(VkPhysicalDevice device) {
            QueueFamilyIndices indices;
            uint32_t queueFamilyCount = 0;
            vkGetPhysicalDeviceQueueFamilyProperties(device, &queueFamilyCount, nullptr);

            std::vector<VkQueueFamilyProperties> queueFamilies(queueFamilyCount);
            vkGetPhysicalDeviceQueueFamilyProperties(device, &queueFamilyCount, queueFamilies.data());

            int i = 0;
            for (const auto& queueFamily : queueFamilies) 
            {
                VkBool32 presentSupport = false;
                vkGetPhysicalDeviceSurfaceSupportKHR(device, i, surface, &presentSupport);
                if (queueFamily.queueFlags & VK_QUEUE_GRAPHICS_BIT) {
                    indices.graphicsFamily = i;
                }
                if (presentSupport)
                    indices.presentFamily = i;
                if (indices.isComplete())
                    break ;
                i++;
            }
            return indices;
        }

        bool checkDeviceExtensionSupport(VkPhysicalDevice tdevice) {
            uint32_t extensionCount;
            vkEnumerateDeviceExtensionProperties(tdevice, nullptr, &extensionCount, nullptr);

            std::vector<VkExtensionProperties> availableExtensions(extensionCount);
            vkEnumerateDeviceExtensionProperties(tdevice, nullptr, &extensionCount, availableExtensions.data());

            std::set<std::string> requiredExtensions(deviceExtensions.begin(), deviceExtensions.end());

            for (const auto& extension : availableExtensions) {
                requiredExtensions.erase(extension.extensionName);
            }

            return requiredExtensions.empty();
        }

        bool isDeviceSuitable(VkPhysicalDevice device) {
            // VkPhysicalDeviceProperties deviceProperties;
            // VkPhysicalDeviceFeatures deviceFeatures;
            // vkGetPhysicalDeviceProperties(device, &deviceProperties);
            // vkGetPhysicalDeviceFeatures(device, &deviceFeatures);

            // return deviceProperties.deviceType == VK_PHYSICAL_DEVICE_TYPE_DISCRETE_GPU &&
            //        deviceFeatures.geometryShader;

            QueueFamilyIndices indices = findQueueFamilies(device);
            bool extensionsSupported = checkDeviceExtensionSupport(device);
            bool swapChainAdequate = false;
            if (extensionsSupported)
            {
                SwapChainSupportDetails swapChainSupport = querySwapChainSupport(device);
                swapChainAdequate = !swapChainSupport.formats.empty() && !swapChainSupport.presentModes.empty();
            }
            return indices.isComplete() && extensionsSupported && swapChainAdequate;
        }

        SwapChainSupportDetails querySwapChainSupport(VkPhysicalDevice device) 
        {
            SwapChainSupportDetails details;
            vkGetPhysicalDeviceSurfaceCapabilitiesKHR(device, surface, &details.capabilities);

            uint32_t formatCount;
            vkGetPhysicalDeviceSurfaceFormatsKHR(device, surface, &formatCount, nullptr);
            if (formatCount != 0)
            {
                details.formats.resize(formatCount);
                vkGetPhysicalDeviceSurfaceFormatsKHR(device, surface, &formatCount, details.formats.data());
            }

            uint32_t presentModeCount;
            vkGetPhysicalDeviceSurfacePresentModesKHR(device, surface, &presentModeCount, nullptr);
            if (presentModeCount != 0)
            {
                details.presentModes.resize(presentModeCount);
                vkGetPhysicalDeviceSurfacePresentModesKHR(device, surface, &presentModeCount, details.presentModes.data());
            }

            return details;
        }

        VkSurfaceFormatKHR chooseSwapSurfaceFormat(const std::vector<VkSurfaceFormatKHR>& availableFormats) 
        {
            for (const auto& availableFormat : availableFormats)
            {
                if (availableFormat.format == VK_FORMAT_B8G8R8A8_SRGB && availableFormat.colorSpace == VK_COLOR_SPACE_SRGB_NONLINEAR_KHR)
                    return availableFormat;
            }
            return availableFormats[0];
        }

        VkPresentModeKHR chooseSwapPresentMode(const std::vector<VkPresentModeKHR>& availablePresentModes)
        {
            for (const auto& availablePresentMode : availablePresentModes)
            {
                if (availablePresentMode == VK_PRESENT_MODE_MAILBOX_KHR)
                    return availablePresentMode;
            }

            return VK_PRESENT_MODE_FIFO_KHR;
        }

        VkExtent2D chooseSwapExtent(const VkSurfaceCapabilitiesKHR& capabilities)
        {
            if (capabilities.currentExtent.width != std::numeric_limits<uint32_t>::max())
                return capabilities.currentExtent;
            else
            {
                int width, height;
                glfwGetFramebufferSize(window, &width, &height);

                VkExtent2D actualExtent = {
                    static_cast<uint32_t>(width),
                    static_cast<uint32_t>(height)
                };

                actualExtent.width = std::clamp(actualExtent.width, capabilities.minImageExtent.width, capabilities.maxImageExtent.width);
                actualExtent.height = std::clamp(actualExtent.height, capabilities.minImageExtent.height, capabilities.maxImageExtent.height);

                return actualExtent;
            }
        }

// init
        void initWindow()
        {
            // if (glfwPlatformSupported(GLFW_PLATFORM_WAYLAND))
            //     glfwInitHint(GLFW_PLATFORM, GLFW_PLATFORM_WAYLAND);
            // else 
            //     std::cout << "No Wayland Support" << std::endl;
            if(!glfwInit())
                throw std::runtime_error("GLFW init failed!");
            if (!glfwVulkanSupported())
                throw std::runtime_error("GLFW reports to have no Vulkan support!");
            glfwWindowHint(GLFW_CLIENT_API, GLFW_NO_API);
            glfwWindowHint(GLFW_RESIZABLE, GLFW_FALSE);
            window = glfwCreateWindow(WIDTH, HEIGHT, "Vulkan", nullptr, nullptr);
        }

        void setupDebugMessenger()
        {
            if (!enableValidationLayers) return;
            VkDebugUtilsMessengerCreateInfoEXT createInfo{};
            populateDebugMessengerCreateInfo(createInfo);
            if (CreateDebugUtilsMessengerEXT(instance, &createInfo, nullptr, &debugMessenger) != VK_SUCCESS) 
                throw std::runtime_error("failed to set up debug messenger!");
        }

        void createSurface() // create via glfw, can be done directly via vk 
        {
            if (glfwCreateWindowSurface(instance, window, nullptr, &surface) != VK_SUCCESS)
                throw std::runtime_error("failed to create window surface!");
        }

        void createInstance()
        {
            if (enableValidationLayers && !checkValidationLayerSupport()) 
                throw std::runtime_error("validation layers requested, but not available!");
    
            VkApplicationInfo appInfo{};
            appInfo.sType = VK_STRUCTURE_TYPE_APPLICATION_INFO;
            appInfo.pApplicationName = "Hello Triangle";
            appInfo.applicationVersion = VK_MAKE_VERSION(1, 0, 0);
            appInfo.pEngineName = "No Engine";
            appInfo.engineVersion = VK_MAKE_VERSION(1, 0, 0);
            appInfo.apiVersion = VK_API_VERSION_1_0;

            auto extensions = getRequiredExtensions();

            VkInstanceCreateInfo createInfo{};
            createInfo.sType = VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO;
            createInfo.pApplicationInfo = &appInfo;
            createInfo.enabledExtensionCount = static_cast<uint32_t>(extensions.size());
            createInfo.ppEnabledExtensionNames = extensions.data();

            VkDebugUtilsMessengerCreateInfoEXT debugCreateInfo{};
            if (enableValidationLayers)
            {
                createInfo.enabledLayerCount = static_cast<uint32_t>(validationLayers.size());
                createInfo.ppEnabledLayerNames = validationLayers.data();
                populateDebugMessengerCreateInfo(debugCreateInfo);
                createInfo.pNext = (VkDebugUtilsMessengerCreateInfoEXT*) &debugCreateInfo;
            }
            else
            {
                createInfo.enabledLayerCount = 0;
                createInfo.pNext = nullptr;
            }

            if (vkCreateInstance(&createInfo, nullptr, &instance) != VK_SUCCESS)
                throw std::runtime_error("failed to create instance!");


        }

        void pickPhysicalDevice()
        {
            uint32_t deviceCount = 0;
            vkEnumeratePhysicalDevices(instance, &deviceCount, nullptr);

            if (deviceCount == 0) 
                throw std::runtime_error("failed to find GPUs with Vulkan support!");
    
            std::vector<VkPhysicalDevice> devices(deviceCount);
            vkEnumeratePhysicalDevices(instance, &deviceCount, devices.data());

            for (const auto& device : devices) {
                if (isDeviceSuitable(device)) {
                    physicalDevice = device;
                    break;
                }
            }

            if (physicalDevice == VK_NULL_HANDLE) 
                throw std::runtime_error("failed to find a suitable GPU!");


        }

        void createLogicalDevice() 
        {
            QueueFamilyIndices                      indices = findQueueFamilies(physicalDevice);
            std::vector<VkDeviceQueueCreateInfo>    queueCreateInfos;
            std::set<uint32_t>                      uniqueQueueFamilies = {indices.graphicsFamily.value(), indices.presentFamily.value()};
            float                                   queuePriority = 1.0f;
    
            for (uint32_t queueFamily : uniqueQueueFamilies) {
                VkDeviceQueueCreateInfo queueCreateInfo{};
                queueCreateInfo.sType = VK_STRUCTURE_TYPE_DEVICE_QUEUE_CREATE_INFO;
                queueCreateInfo.queueFamilyIndex = queueFamily;
                queueCreateInfo.queueCount = 1;
                queueCreateInfo.pQueuePriorities = &queuePriority;
                queueCreateInfos.push_back(queueCreateInfo);
            }


            VkPhysicalDeviceFeatures    deviceFeatures{};
            VkDeviceCreateInfo          createInfo{};

            createInfo.sType = VK_STRUCTURE_TYPE_DEVICE_CREATE_INFO;
            createInfo.queueCreateInfoCount = static_cast<uint32_t>(queueCreateInfos.size());
            createInfo.pQueueCreateInfos = queueCreateInfos.data();
            createInfo.pEnabledFeatures = &deviceFeatures;
            createInfo.enabledExtensionCount = static_cast<uint32_t>(deviceExtensions.size());
            createInfo.ppEnabledExtensionNames = deviceExtensions.data();

            if (enableValidationLayers) // ignored by up to date impl 
            { 
                createInfo.enabledLayerCount = static_cast<uint32_t>(validationLayers.size());
                createInfo.ppEnabledLayerNames = validationLayers.data();
            } 
            else
                createInfo.enabledLayerCount = 0;

            if (vkCreateDevice(physicalDevice, &createInfo, nullptr, &device) != VK_SUCCESS)
                throw std::runtime_error("failed to create logical device!");

            vkGetDeviceQueue(device, indices.graphicsFamily.value(), 0, &graphicsQueue);
            vkGetDeviceQueue(device, indices.presentFamily.value(), 0, &presentQueue);

        }

        void createSwapChain()
        {
            SwapChainSupportDetails swapChainSupport = querySwapChainSupport(physicalDevice);

            VkSurfaceFormatKHR      surfaceFormat = chooseSwapSurfaceFormat(swapChainSupport.formats);
            VkPresentModeKHR        presentMode = chooseSwapPresentMode(swapChainSupport.presentModes);
            VkExtent2D              extent = chooseSwapExtent(swapChainSupport.capabilities);

            uint32_t                imageCount = swapChainSupport.capabilities.minImageCount + 1;
            if (swapChainSupport.capabilities.maxImageCount > 0 && imageCount > swapChainSupport.capabilities.maxImageCount) 
                imageCount = swapChainSupport.capabilities.maxImageCount;

            VkSwapchainCreateInfoKHR createInfo{};
            createInfo.sType = VK_STRUCTURE_TYPE_SWAPCHAIN_CREATE_INFO_KHR;
            createInfo.surface = surface;
            createInfo.minImageCount = imageCount;
            createInfo.imageFormat = surfaceFormat.format;
            createInfo.imageColorSpace = surfaceFormat.colorSpace;
            createInfo.imageExtent = extent;
            createInfo.imageArrayLayers = 1;
            createInfo.imageUsage = VK_IMAGE_USAGE_COLOR_ATTACHMENT_BIT;

            QueueFamilyIndices indices = findQueueFamilies(physicalDevice);
            uint32_t queueFamilyIndices[] = {indices.graphicsFamily.value(), indices.presentFamily.value()};

            if (indices.graphicsFamily != indices.presentFamily)
            {
                std::cout << "concurrent mode for queue" << std::endl;
                createInfo.imageSharingMode = VK_SHARING_MODE_CONCURRENT;
                createInfo.queueFamilyIndexCount = 2;
                createInfo.pQueueFamilyIndices = queueFamilyIndices;
            }
            else
            {
                createInfo.imageSharingMode = VK_SHARING_MODE_EXCLUSIVE;
                createInfo.queueFamilyIndexCount = 0; // Optional
                createInfo.pQueueFamilyIndices = nullptr; // Optional
            }

            createInfo.preTransform = swapChainSupport.capabilities.currentTransform;
            createInfo.compositeAlpha = VK_COMPOSITE_ALPHA_OPAQUE_BIT_KHR;
            createInfo.presentMode = presentMode;
            createInfo.clipped = VK_TRUE;
            createInfo.oldSwapchain = VK_NULL_HANDLE;

            VkResult result = vkCreateSwapchainKHR(device, &createInfo, nullptr, &swapChain);
            if (result != VK_SUCCESS)
            {
                std::string res_str(string_VkResult(result));
                throw std::runtime_error("failed to create swap chain! error: " + res_str);
            }

            vkGetSwapchainImagesKHR(device, swapChain, &imageCount, nullptr);
            swapChainImages.resize(imageCount);
            vkGetSwapchainImagesKHR(device, swapChain, &imageCount, swapChainImages.data());

            swapChainImageFormat = surfaceFormat.format;
            swapChainExtent = extent;

        }

        void createImageViews()
        {
            swapChainImageViews.resize(swapChainImages.size());

            for (size_t i = 0; i < swapChainImages.size(); i++) 
            {
                VkImageViewCreateInfo createInfo{};
                createInfo.sType = VK_STRUCTURE_TYPE_IMAGE_VIEW_CREATE_INFO;
                createInfo.image = swapChainImages[i];
                createInfo.viewType = VK_IMAGE_VIEW_TYPE_2D;
                createInfo.format = swapChainImageFormat;
                createInfo.components.r = VK_COMPONENT_SWIZZLE_IDENTITY;
                createInfo.components.g = VK_COMPONENT_SWIZZLE_IDENTITY;
                createInfo.components.b = VK_COMPONENT_SWIZZLE_IDENTITY;
                createInfo.components.a = VK_COMPONENT_SWIZZLE_IDENTITY;
                createInfo.subresourceRange.aspectMask = VK_IMAGE_ASPECT_COLOR_BIT;
                createInfo.subresourceRange.baseMipLevel = 0;
                createInfo.subresourceRange.levelCount = 1;
                createInfo.subresourceRange.baseArrayLayer = 0;
                createInfo.subresourceRange.layerCount = 1;
                if (vkCreateImageView(device, &createInfo, nullptr, &swapChainImageViews[i]) != VK_SUCCESS)
                    throw std::runtime_error("échec de la création d'une image view!");
            }

        }

        VkShaderModule createShaderModule(const std::vector<char>& code) 
        {
            VkShaderModuleCreateInfo createInfo{};

            createInfo.sType = VK_STRUCTURE_TYPE_SHADER_MODULE_CREATE_INFO;
            createInfo.codeSize = code.size();
            createInfo.pCode = reinterpret_cast<const uint32_t*>(code.data());

            VkShaderModule shaderModule;
            if (vkCreateShaderModule(device, &createInfo, nullptr, &shaderModule) != VK_SUCCESS)
                throw std::runtime_error("failed to create shader module!");

            return shaderModule;
        }

        //  to tell vulkan about the framebuffer attachment : color and depth buffer
        void createRenderPass()
        {
            // we use 1 color buffer which format match a image 
            VkAttachmentDescription colorAttachment{};
            colorAttachment.format = swapChainImageFormat;
            colorAttachment.samples = VK_SAMPLE_COUNT_1_BIT;
            // determine what to do with data attachment before and after rendering
            /*
            for loadOp:
                VK_ATTACHMENT_LOAD_OP_LOAD: Preserve the existing contents of the attachment
                VK_ATTACHMENT_LOAD_OP_CLEAR: Clear the values to a constant at the start
                VK_ATTACHMENT_LOAD_OP_DONT_CARE: Existing contents are undefined; we don't care about them
            */
            colorAttachment.loadOp = VK_ATTACHMENT_LOAD_OP_CLEAR;
            /*
            for storeOp:
                VK_ATTACHMENT_STORE_OP_STORE: Rendered contents will be stored in memory and can be read later
                VK_ATTACHMENT_STORE_OP_DONT_CARE: Contents of the framebuffer will be undefined after the rendering operation
            */
            colorAttachment.storeOp = VK_ATTACHMENT_STORE_OP_STORE;
            // idem for stencil buffer
            // but we dont have one so dont care
            colorAttachment.stencilLoadOp = VK_ATTACHMENT_LOAD_OP_DONT_CARE;
            colorAttachment.stencilStoreOp = VK_ATTACHMENT_STORE_OP_DONT_CARE;

            // layout of pixel in memory
            /*
            can be (other exist):
                VK_IMAGE_LAYOUT_COLOR_ATTACHMENT_OPTIMAL: Images used as color attachment
                VK_IMAGE_LAYOUT_PRESENT_SRC_KHR: Images to be presented in the swap chain
                VK_IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL: Images to be used as destination for a memory copy operation
            */
            colorAttachment.initialLayout = VK_IMAGE_LAYOUT_UNDEFINED;
            colorAttachment.finalLayout = VK_IMAGE_LAYOUT_PRESENT_SRC_KHR;

            // one render pass can combine multiple subpass (subsequent rendering operation dependent of previous framebuffer content)
            // example : post-process effect. Grouping into one render pass alllow optimization from vulkan
            // we only use one here

            // we pass to the subpass an attachment reference
            VkAttachmentReference colorAttachmentRef{};
            // index of attachemtn in array 
            colorAttachmentRef.attachment = 0;
            colorAttachmentRef.layout = VK_IMAGE_LAYOUT_COLOR_ATTACHMENT_OPTIMAL;

            // subpass : graphic, with 1 attachment ref
            VkSubpassDescription subpass{};
            subpass.pipelineBindPoint = VK_PIPELINE_BIND_POINT_GRAPHICS;
            subpass.colorAttachmentCount = 1;
            subpass.pColorAttachments = &colorAttachmentRef;
            // our attachment here is referenced from the fragment shader wit hthe directive :layout(location = 0) out vec4 outColor 
            //  subpass can hold several type of attachment
            /*
                pInputAttachments: Attachments that are read from a shader
                pResolveAttachments: Attachments used for multisampling color attachments
                pDepthStencilAttachment: Attachment for depth and stencil data
                pPreserveAttachments: Attachments that are not used by this subpass, but for which the data must be preserved
            */


            // renderpass CI
            // we pass aray of attachment and array of subpass
            VkRenderPassCreateInfo renderPassInfo{};
            renderPassInfo.sType = VK_STRUCTURE_TYPE_RENDER_PASS_CREATE_INFO;
            renderPassInfo.attachmentCount = 1;
            renderPassInfo.pAttachments = &colorAttachment;
            renderPassInfo.subpassCount = 1;
            renderPassInfo.pSubpasses = &subpass;

            if (vkCreateRenderPass(device, &renderPassInfo, nullptr, &renderPass) != VK_SUCCESS)
                throw std::runtime_error("failed to create render pass!");

        }

        // render process
        /*
            Shader stages: the shader modules that define the functionality of the programmable stages of the graphics pipeline
            Fixed-function state: all of the structures that define the fixed-function stages of the pipeline, like input assembly, rasterizer, viewport and color blending
            Pipeline layout: the uniform and push values referenced by the shader that can be updated at draw time
            Render pass: the attachments referenced by the pipeline stages and their usage
        */
        void createGraphicsPipeline() 
        {
            // load shader binary code
            auto vertShaderCode = readFile("shaders/vert.spv");
            auto fragShaderCode = readFile("shaders/frag.spv");

            // shader module from shader bin code
            VkShaderModule vertShaderModule = createShaderModule(vertShaderCode);
            VkShaderModule fragShaderModule = createShaderModule(fragShaderCode);

            // shader CI for vertex shader
            VkPipelineShaderStageCreateInfo vertShaderStageInfo{};
            vertShaderStageInfo.sType = VK_STRUCTURE_TYPE_PIPELINE_SHADER_STAGE_CREATE_INFO;
            vertShaderStageInfo.stage = VK_SHADER_STAGE_VERTEX_BIT;
            vertShaderStageInfo.module = vertShaderModule;
            vertShaderStageInfo.pName = "main";

            // shader CI for fragemnt shader
            VkPipelineShaderStageCreateInfo fragShaderStageInfo{};
            fragShaderStageInfo.sType = VK_STRUCTURE_TYPE_PIPELINE_SHADER_STAGE_CREATE_INFO;
            fragShaderStageInfo.stage = VK_SHADER_STAGE_FRAGMENT_BIT;
            fragShaderStageInfo.module = fragShaderModule;
            fragShaderStageInfo.pName = "main";


            // link shader CI as array
            VkPipelineShaderStageCreateInfo shaderStages[] = {vertShaderStageInfo, fragShaderStageInfo};

            // for dynamique change (only viewport and few things)
            // we MUST specify valu at runtime if we go dynamic
            // list of what we want as dynamique
            std::vector<VkDynamicState> dynamicStates = {
                VK_DYNAMIC_STATE_VIEWPORT,
                VK_DYNAMIC_STATE_SCISSOR
            };

            // dynamicstate CI
            VkPipelineDynamicStateCreateInfo dynamicState{};
            dynamicState.sType = VK_STRUCTURE_TYPE_PIPELINE_DYNAMIC_STATE_CREATE_INFO;
            dynamicState.dynamicStateCount = static_cast<uint32_t>(dynamicStates.size());
            dynamicState.pDynamicStates = dynamicStates.data();

            // vertex input CI (describe format vertex give nto vertex shader)
            // filled for no data to load
            VkPipelineVertexInputStateCreateInfo vertexInputInfo{};
            vertexInputInfo.sType = VK_STRUCTURE_TYPE_PIPELINE_VERTEX_INPUT_STATE_CREATE_INFO;
            vertexInputInfo.vertexBindingDescriptionCount = 0;
            vertexInputInfo.pVertexBindingDescriptions = nullptr; // Optional
            vertexInputInfo.vertexAttributeDescriptionCount = 0;
            vertexInputInfo.pVertexAttributeDescriptions = nullptr; // Optional

            // input assembly CI (describe how to use input vertice: draw pont, line, triangle... reuse vertices,..)
            VkPipelineInputAssemblyStateCreateInfo inputAssembly{};
            inputAssembly.sType = VK_STRUCTURE_TYPE_PIPELINE_INPUT_ASSEMBLY_STATE_CREATE_INFO;
            inputAssembly.topology = VK_PRIMITIVE_TOPOLOGY_TRIANGLE_LIST;
            inputAssembly.primitiveRestartEnable = VK_FALSE;

            // (D) viewport
            // region of framebuffer that the output will be rendered to
            // VkViewport viewport{};
            // viewport.x = 0.0f;
            // viewport.y = 0.0f;
            // viewport.width = (float) swapChainExtent.width;
            // viewport.height = (float) swapChainExtent.height;
            // viewport.minDepth = 0.0f;
            // viewport.maxDepth = 1.0f;

            // (D) scissor 
            // region of framebuffer where pixel are allowed to be stored (outisde its discarded)
            // VkRect2D scissor{};
            // scissor.offset = {0, 0};
            // scissor.extent = swapChainExtent;

            // if viewport/scissor dynamic, just need to set their count in CI
            VkPipelineViewportStateCreateInfo viewportState{};
            viewportState.sType = VK_STRUCTURE_TYPE_PIPELINE_VIEWPORT_STATE_CREATE_INFO;
            viewportState.viewportCount = 1;
            viewportState.scissorCount = 1;


            // else we load them 
            // VkPipelineViewportStateCreateInfo viewportState{};
            // viewportState.sType = VK_STRUCTURE_TYPE_PIPELINE_VIEWPORT_STATE_CREATE_INFO;
            // viewportState.viewportCount = 1;
            // viewportState.pViewports = &viewport;
            // viewportState.scissorCount = 1;
            // viewportState.pScissors = &scissor;

            // (G) we can set multpitple viewport/scissor in one pipeline, need to be enabled at (logical) device level


            // rasterizer (take geometry and fill with fragment)
            // perform depth testing, face culling, scissor test
            VkPipelineRasterizationStateCreateInfo rasterizer{};
            rasterizer.sType = VK_STRUCTURE_TYPE_PIPELINE_RASTERIZATION_STATE_CREATE_INFO;

            // (G) fragment that are beyon near and far plane are clamped instedd of discarded (usefull for shadowmap) 
            // not enabled here
            rasterizer.depthClampEnable = VK_FALSE;

            // enable/disable rasterization (disable any output to framebuffer)
            rasterizer.rasterizerDiscardEnable = VK_FALSE;

            // (G) determin how fragment are generated, (G) for other than FILL
            /*
                VK_POLYGON_MODE_FILL: fill the area of the polygon with fragments
                VK_POLYGON_MODE_LINE: polygon edges are drawn as lines
                VK_POLYGON_MODE_POINT: polygon vertices are drawn as point
            */
            rasterizer.polygonMode = VK_POLYGON_MODE_FILL;

            // thickness of line in count of fragment, maximum depend of hardware, > 1.0f need (G) wideLine
            rasterizer.lineWidth = 1.0f;

            // type of face culling and vertex order 
            rasterizer.cullMode = VK_CULL_MODE_BACK_BIT;
            rasterizer.frontFace = VK_FRONT_FACE_CLOCKWISE;

            // alter depth values
            // not enabled here
            rasterizer.depthBiasEnable = VK_FALSE;
            rasterizer.depthBiasConstantFactor = 0.0f; // Optional
            rasterizer.depthBiasClamp = 0.0f; // Optional
            rasterizer.depthBiasSlopeFactor = 0.0f; // Optional


            // (G) enable multisample anti-aliasing, a way to perform anti aliasing 
            // not enabled here
            VkPipelineMultisampleStateCreateInfo multisampling{};
            multisampling.sType = VK_STRUCTURE_TYPE_PIPELINE_MULTISAMPLE_STATE_CREATE_INFO;
            multisampling.sampleShadingEnable = VK_FALSE;
            multisampling.rasterizationSamples = VK_SAMPLE_COUNT_1_BIT;
            multisampling.minSampleShading = 1.0f; // Optional
            multisampling.pSampleMask = nullptr; // Optional
            multisampling.alphaToCoverageEnable = VK_FALSE; // Optional
            multisampling.alphaToOneEnable = VK_FALSE; // Optional

            // we could configure a depth/stencil test if we re using a depht/stencil buffer
            // we dont do that now
            // VkPipelineDepthStencilStateCreateInfo 


            // color blending (combine color given by fragment shader and color of framebuffer)
            // can be an array for different framebuffer
            // we dont do that 
            VkPipelineColorBlendAttachmentState colorBlendAttachment{};
            colorBlendAttachment.colorWriteMask = VK_COLOR_COMPONENT_R_BIT | VK_COLOR_COMPONENT_G_BIT | VK_COLOR_COMPONENT_B_BIT | VK_COLOR_COMPONENT_A_BIT;
            colorBlendAttachment.blendEnable = VK_FALSE;
            colorBlendAttachment.srcColorBlendFactor = VK_BLEND_FACTOR_ONE; // Optional
            colorBlendAttachment.dstColorBlendFactor = VK_BLEND_FACTOR_ZERO; // Optional
            colorBlendAttachment.colorBlendOp = VK_BLEND_OP_ADD; // Optional
            colorBlendAttachment.srcAlphaBlendFactor = VK_BLEND_FACTOR_ONE; // Optional
            colorBlendAttachment.dstAlphaBlendFactor = VK_BLEND_FACTOR_ZERO; // Optional
            colorBlendAttachment.alphaBlendOp = VK_BLEND_OP_ADD; // Optional

            // color blending CI
            VkPipelineColorBlendStateCreateInfo colorBlending{};
            colorBlending.sType = VK_STRUCTURE_TYPE_PIPELINE_COLOR_BLEND_STATE_CREATE_INFO;
            colorBlending.logicOpEnable = VK_FALSE;
            colorBlending.logicOp = VK_LOGIC_OP_COPY; // Optional
            colorBlending.attachmentCount = 1;
            colorBlending.pAttachments = &colorBlendAttachment;
            colorBlending.blendConstants[0] = 0.0f; // Optional
            colorBlending.blendConstants[1] = 0.0f; // Optional
            colorBlending.blendConstants[2] = 0.0f; // Optional
            colorBlending.blendConstants[3] = 0.0f; // Optional


            // we can use `uniform` values in shaders, basically global similar to dynamic state values, 
            // which can be changed a runtime
            // those a defined in a VkPipelineLayout object, defined as a class member
            // we dont do that
            // but we need a empty object to create the pipeline
            
            // pipeline layout CI
            VkPipelineLayoutCreateInfo pipelineLayoutInfo{};
            pipelineLayoutInfo.sType = VK_STRUCTURE_TYPE_PIPELINE_LAYOUT_CREATE_INFO;
            pipelineLayoutInfo.setLayoutCount = 0; // Optional
            pipelineLayoutInfo.pSetLayouts = nullptr; // Optional
            pipelineLayoutInfo.pushConstantRangeCount = 0; // Optional
            pipelineLayoutInfo.pPushConstantRanges = nullptr; // Optional

            // we cant also use `push constant` to pass value to shader, but well see that another time
            if (vkCreatePipelineLayout(device, &pipelineLayoutInfo, nullptr, &pipelineLayout) != VK_SUCCESS) 
                throw std::runtime_error("failed to create pipeline layout!");


            //  yeepee! creating the pipeline ! at last!
            VkGraphicsPipelineCreateInfo pipelineInfo{};
            pipelineInfo.sType = VK_STRUCTURE_TYPE_GRAPHICS_PIPELINE_CREATE_INFO;
            // referencing shaders stages
            pipelineInfo.stageCount = 2;
            pipelineInfo.pStages = shaderStages;
            // other stage and fxed function
            pipelineInfo.pVertexInputState = &vertexInputInfo;
            pipelineInfo.pInputAssemblyState = &inputAssembly;
            pipelineInfo.pViewportState = &viewportState;
            pipelineInfo.pRasterizationState = &rasterizer;
            pipelineInfo.pMultisampleState = &multisampling;
            pipelineInfo.pDepthStencilState = nullptr; // Optional
            pipelineInfo.pColorBlendState = &colorBlending;
            pipelineInfo.pDynamicState = &dynamicState;
            // then layout
            pipelineInfo.layout = pipelineLayout;
            // then render pass wit hindex of subpass
            pipelineInfo.renderPass = renderPass;
            pipelineInfo.subpass = 0;

            // opttionnal, allow to create pipeline derivate from another parent pipelein
            // only used if VK_PIPELINE_CREATE_DERIVATIVE_BIT flag is set
            pipelineInfo.basePipelineHandle = VK_NULL_HANDLE; // Optional
            pipelineInfo.basePipelineIndex = -1; // Optional

            // create pipeline, can create multpiple pipeline at once and used cached data to create pipelien faster
            if (vkCreateGraphicsPipelines(device, VK_NULL_HANDLE, 1, &pipelineInfo, nullptr, &graphicsPipeline) != VK_SUCCESS)
                throw std::runtime_error("failed to create graphics pipeline!");

            //  destroy shader module, as they are now loaded in the pipeline
            vkDestroyShaderModule(device, fragShaderModule, nullptr);
            vkDestroyShaderModule(device, vertShaderModule, nullptr);


        }

        // framebuffer object ???
        // bound to one or more image via imageview
        void createFramebuffers()
        {
            swapChainFramebuffers.resize(swapChainImageViews.size());
            for (size_t i = 0; i < swapChainImageViews.size(); i++)
            {
                VkImageView attachments[] = {
                    swapChainImageViews[i]
                };

                VkFramebufferCreateInfo framebufferInfo{};
                framebufferInfo.sType = VK_STRUCTURE_TYPE_FRAMEBUFFER_CREATE_INFO;
                framebufferInfo.renderPass = renderPass;
                framebufferInfo.attachmentCount = 1;
                framebufferInfo.pAttachments = attachments;
                framebufferInfo.width = swapChainExtent.width;
                framebufferInfo.height = swapChainExtent.height;
                framebufferInfo.layers = 1;

                if (vkCreateFramebuffer(device, &framebufferInfo, nullptr, &swapChainFramebuffers[i]) != VK_SUCCESS)
                    throw std::runtime_error("failed to create framebuffer!");
            }
        }

        // command pool for submit command to vulkan
        // one command pool for one family of queue
        void createCommandPool()
        {
            QueueFamilyIndices queueFamilyIndices = findQueueFamilies(physicalDevice);

            VkCommandPoolCreateInfo poolInfo{};
            poolInfo.sType = VK_STRUCTURE_TYPE_COMMAND_POOL_CREATE_INFO;
            // allow to rerecord command buffer individually
            poolInfo.flags = VK_COMMAND_POOL_CREATE_RESET_COMMAND_BUFFER_BIT;
            poolInfo.queueFamilyIndex = queueFamilyIndices.graphicsFamily.value();

            if (vkCreateCommandPool(device, &poolInfo, nullptr, &commandPool) != VK_SUCCESS)
                throw std::runtime_error("failed to create command pool!");
        }

        // create command buffer
        // command buffer are automatically freed whe ncommand is allocated
        void    createCommandBuffer() 
        {
            VkCommandBufferAllocateInfo allocInfo{};
            allocInfo.sType = VK_STRUCTURE_TYPE_COMMAND_BUFFER_ALLOCATE_INFO;
            allocInfo.commandPool = commandPool;
            /*
                VK_COMMAND_BUFFER_LEVEL_PRIMARY: Can be submitted to a queue for execution, but cannot be called from other command buffers.
                VK_COMMAND_BUFFER_LEVEL_SECONDARY: Cannot be submitted directly, but can be called from primary command buffers.
            */
            allocInfo.level = VK_COMMAND_BUFFER_LEVEL_PRIMARY;
            allocInfo.commandBufferCount = 1;

            if (vkAllocateCommandBuffers(device, &allocInfo, &commandBuffer) != VK_SUCCESS)
                throw std::runtime_error("failed to allocate command buffers!");

        }

        void recordCommandBuffer(VkCommandBuffer commandBuffer, uint32_t imageIndex)
        {
            VkCommandBufferBeginInfo beginInfo{};
            beginInfo.sType = VK_STRUCTURE_TYPE_COMMAND_BUFFER_BEGIN_INFO;
            /* 
                VK_COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT: The command buffer will be rerecorded right after executing it once.
                VK_COMMAND_BUFFER_USAGE_RENDER_PASS_CONTINUE_BIT: This is a secondary command buffer that will be entirely within a single render pass.
                VK_COMMAND_BUFFER_USAGE_SIMULTANEOUS_USE_BIT: The command buffer can be resubmitted while it is also already pending execution
            */
            beginInfo.flags = 0; // Optional
            beginInfo.pInheritanceInfo = nullptr; // Optional, only for secondary command
            
            // if command already recorded, this will reset it, we cant append command to buffer
            if (vkBeginCommandBuffer(commandBuffer, &beginInfo) != VK_SUCCESS)
                throw std::runtime_error("failed to begin recording command buffer!");

            // we start drawing by beginning the render pass 
            VkRenderPassBeginInfo renderPassInfo{};
            renderPassInfo.sType = VK_STRUCTURE_TYPE_RENDER_PASS_BEGIN_INFO;
            // bind renderpass and framebuffer to use
            renderPassInfo.renderPass = renderPass;
            renderPassInfo.framebuffer = swapChainFramebuffers[imageIndex];
            // size of render area, anything outise will be undefined
            renderPassInfo.renderArea.offset = {0, 0};
            renderPassInfo.renderArea.extent = swapChainExtent;


            // some colors idk, for VK_ATTACHMENT_LOAD_OP_CLEAR (we used it i nrender pass?)
            VkClearValue clearColor = {{{0.0f, 0.0f, 0.0f, 1.0f}}};
            renderPassInfo.clearValueCount = 1;
            renderPassInfo.pClearValues = &clearColor;

            // we wnat tot do a renderpass
            vkCmdBeginRenderPass(commandBuffer, &renderPassInfo, VK_SUBPASS_CONTENTS_INLINE);
            // we bind the pipeline
            vkCmdBindPipeline(commandBuffer, VK_PIPELINE_BIND_POINT_GRAPHICS, graphicsPipeline);

            // we say we want some (D) dynamic stuff, here we go :
            // we need to define them now
            VkViewport viewport{};
            viewport.x = 0.0f;
            viewport.y = 0.0f;
            viewport.width = static_cast<float>(swapChainExtent.width);
            viewport.height = static_cast<float>(swapChainExtent.height);
            viewport.minDepth = 0.0f;
            viewport.maxDepth = 1.0f;
            vkCmdSetViewport(commandBuffer, 0, 1, &viewport);

            VkRect2D scissor{};
            scissor.offset = {0, 0};
            scissor.extent = swapChainExtent;
            vkCmdSetScissor(commandBuffer, 0, 1, &scissor);

            // now we can create the drw command
            // very simple because we alread ystuffed a lot of data inside cmdbuffer
            // paramter: cmd, number of vertex to draw, instaced rendering (we dont do that), vertex offset in vertex buffer, same for isntace 
            vkCmdDraw(commandBuffer, 3, 1, 0, 0);
            // end renderpass
            vkCmdEndRenderPass(commandBuffer);
            // end cmdbuffer
            if (vkEndCommandBuffer(commandBuffer) != VK_SUCCESS) 
                throw std::runtime_error("failed to record command buffer!");

        }


        void initVulkan() 
        {
            createInstance();
            setupDebugMessenger();
            createSurface();
            pickPhysicalDevice();
            createLogicalDevice();
            createSwapChain();
            createImageViews();
            createRenderPass();
            createGraphicsPipeline();
            createFramebuffers();
            createCommandPool();
            createCommandBuffer();
        }


        void mainLoop() 
        {
            while (!glfwWindowShouldClose(window)) {
                    glfwPollEvents();
                }
        }

        void cleanup() 
        {
            vkDestroyCommandPool(device, commandPool, nullptr);


            for (auto framebuffer : swapChainFramebuffers) {
                vkDestroyFramebuffer(device, framebuffer, nullptr);
            }

            vkDestroyPipeline(device, graphicsPipeline, nullptr);
            vkDestroyPipelineLayout(device, pipelineLayout, nullptr);
            vkDestroyRenderPass(device, renderPass, nullptr);
            for (auto imageView : swapChainImageViews) {
                vkDestroyImageView(device, imageView, nullptr);
            }
            
            vkDestroySwapchainKHR(device, swapChain, nullptr);
            vkDestroyDevice(device, nullptr);

            if (enableValidationLayers) 
                DestroyDebugUtilsMessengerEXT(instance, debugMessenger, nullptr);
            vkDestroySurfaceKHR(instance, surface, nullptr);
            vkDestroyInstance(instance, nullptr);

            glfwDestroyWindow(window);

            glfwTerminate();
        }
};


int main() {
    TApp app;

    try {
        app.run();
    } catch (const std::exception& e) {
        std::cerr << e.what() << std::endl;
        return EXIT_FAILURE;
    }

    return EXIT_SUCCESS;
}