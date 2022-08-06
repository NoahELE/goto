import {
  Button,
  Center,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Heading,
  HStack,
  Input,
  Link,
  VStack,
} from '@chakra-ui/react'
import { useState } from 'react'
import { SubmitHandler, useForm } from 'react-hook-form'

type UrlForm = {
  url: string
}

export default function App() {
  const [link, setLink] = useState('')

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<UrlForm>()

  const onSubmit: SubmitHandler<UrlForm> = async function (values) {
    const r = await fetch('/add', {
      method: 'post',
      body: JSON.stringify(values),
    })
    const json = await r.json()
    setLink(window.location.origin + '/' + json.key)
  }

  return (
    <>
      <Heading
        color="teal.500"
        as="h1"
        m={50}
        size="3xl"
        fontStyle="italic"
        fontWeight="black"
        textAlign="center"
      >
        GOTO
      </Heading>
      <Center m={20}>
        <VStack>
          <form onSubmit={handleSubmit(onSubmit)}>
            <FormControl isInvalid={errors.url !== undefined} w={700}>
              <HStack>
                <FormLabel
                  htmlFor="url-input"
                  fontSize="3xl"
                  fontWeight="bold"
                  color="gray.600"
                >
                  URL:
                </FormLabel>
                <Input
                  id="url-input"
                  placeholder="please input a complete url"
                  size="lg"
                  focusBorderColor="teal.500"
                  {...register('url', {
                    required: 'a url is required',
                    maxLength: {
                      value: 128,
                      message: 'url is too long',
                    },
                    validate: isValidUrl,
                  })}
                />
              </HStack>
              <FormErrorMessage fontSize="lg" fontWeight="bold">
                {errors.url?.message}
              </FormErrorMessage>
            </FormControl>
            <Center m={10}>
              <Button colorScheme="teal" isLoading={isSubmitting} type="submit">
                Generate Short URL
              </Button>
            </Center>
          </form>
          <br />
          <Link href={link} isExternal>
            {link}
          </Link>
        </VStack>
      </Center>
    </>
  )
}

function isValidUrl(s: string) {
  try {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const _ = new URL(s)
  } catch (err) {
    return 'not a valid url'
  }
  return true
}
